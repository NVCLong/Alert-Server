package controller

import (
	"time"

	"github.com/NVCLong/Alert-Server/modules/middleware"
	workflowController "github.com/NVCLong/Alert-Server/modules/work-flow/controller"
	redisServices "github.com/NVCLong/Alert-Server/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(timeout time.Duration, db *gorm.DB, gin *gin.Engine, cacheService redisServices.AbstractCacheService) {

	adminMiddleware := middleware.NewAdminMiddleware(db, cacheService)
	// publicRouter := gin.Group("")
	// Register route for other controller

	protectedRouter := gin.Group("/api")
	//register route related to admin function
	protectedRouter.Use(adminMiddleware.GetAdminHandlerFunc())

	workflowController.NewWorkFlowController(timeout, db, protectedRouter)

}
