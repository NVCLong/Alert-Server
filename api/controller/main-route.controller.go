package controller

import (
	"time"

	"github.com/NVCLong/Alert-Server/api/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(timeout time.Duration, db *gorm.DB, gin *gin.Engine) {
	// publicRouter := gin.Group("")
	// Register route for other controller

	protectedRouter := gin.Group("")
	//register route related to admin function
	protectedRouter.Use(middleware.AdminMiddleware(db))

	NewWorkFlowController(timeout, db, protectedRouter)

}
