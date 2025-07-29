package main

import (
	"net/http"
	"os"
	"strconv"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDBStr := os.Getenv("REDIS_DB")
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		redisDB = 0
	}

	rClient := NewRedisClient(redisAddr, redisPassword, redisDB)

	router := gin.Default()

	router.GET("/clipboard/:key/:value", func(c *gin.Context) {
		key := c.Param("key")
		value := c.Param("value")

		if err := rClient.Set(key, value); err != nil {
			fmt.Printf("Error : %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set value"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Value set successfully"})
	})

	router.GET("/clipboard/:key", func(c *gin.Context) {
		key := c.Param("key")

		val, err := rClient.Get(key)
		if err != nil {
			if err == RedisNil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Value not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get value"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"key": key, "value": val})
	})

	router.Run(":8080") // Start server on port 8080
}
