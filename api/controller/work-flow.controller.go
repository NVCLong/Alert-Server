package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewWorkFlowController(timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	workflowGroup := group.Group("/work-flow")

	workflowGroup.GET("/all", getAllWorkFlows)
}

func getAllWorkFlows(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Getting all work flow"})
}
