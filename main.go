package main

import (
	"github.com/NVCLong/Alert-Server/database"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	database.ConnectDatabase()
	router.Run("localhost:8080")

}
