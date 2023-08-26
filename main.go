package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var client *redis.Client

func init() {
	// Configuración de conexión a Redis
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Cambia esto según la configuración de tu servidor Redis
		Password: "",               // Contraseña (si está configurada)
		DB:       0,                // Número de base de datos
	})
}

func main() {
	// Cerrar la conexión al final
	defer client.Close()

	r := gin.Default()

	r.POST("/set", setValue)
	r.GET("/get/:key", getValue)

	r.Run(":8080")
}

func setValue(c *gin.Context) {
	var requestData struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	err := client.Set(ctx, requestData.Key, requestData.Value, 10*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al almacenar en caché"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Datos almacenados en caché correctamente"})
}

func getValue(c *gin.Context) {
	key := c.Param("key")

	ctx := context.Background()
	value, err := client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Clave no encontrada en caché"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener datos de caché"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": value})
}
