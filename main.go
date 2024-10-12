package main

import (
	"fmt"
	"time"

	boostrap "github.com/NVCLong/Alert-Server/bootstrap"
	"github.com/NVCLong/Alert-Server/database"
	mainController "github.com/NVCLong/Alert-Server/modules"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	db := database.ConnectDatabase()

	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "UserId"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))
	timeout := time.Duration(30) * time.Second
	router.Group("/api")
	mainController.Setup(timeout, db, router)
	router.Run(fmt.Sprintf("localhost:%s", boostrap.GetEnv(boostrap.EnvServerPort)))
}
