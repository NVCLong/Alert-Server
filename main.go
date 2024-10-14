package main

import (
	"context"
	"time"

	boostrap "github.com/NVCLong/Alert-Server/bootstrap"
	"github.com/NVCLong/Alert-Server/database"
	mainController "github.com/NVCLong/Alert-Server/modules"
	"github.com/NVCLong/Alert-Server/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	db := database.ConnectDatabase()
	port := boostrap.GetEnv(boostrap.EnvServerPort)
	redisClient := redis.NewRedisConnection()
	cacheService := redis.NewCacheService(*redisClient, context.Background())
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "UserId"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the API!",
		})
	})
	timeout := time.Duration(30) * time.Second
	router.Group("/api")
	mainController.Setup(timeout, db, router, cacheService)
	router.Run(":" + port)
}
