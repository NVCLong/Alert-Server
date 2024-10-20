package controller

import (
	"log"
	"net/http"
	"time"

	abstractrepo "github.com/NVCLong/Alert-Server/database/abstract-repo"
	"github.com/NVCLong/Alert-Server/dto"
	"github.com/NVCLong/Alert-Server/modules/work-flow/repository"
	"github.com/NVCLong/Alert-Server/modules/work-flow/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewWorkFlowController(timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	workflowGroup := group.Group("/work-flow")

	workFlowRepository := repository.NewWorkFlowRepository(db, abstractrepo.WorkFlowTable)

	workFlowService := service.NewWorkFlowService(workFlowRepository)

	workflowGroup.GET("/all", func(ctx *gin.Context) {
		getAllWorkFlows(ctx, workFlowService)
	})
	workflowGroup.POST("/create", func(ctx *gin.Context) {
		createWorkFlow(ctx, workFlowService)
	})
}

func getAllWorkFlows(c *gin.Context, workFlowService service.WorkFlowAbstractService) {
	workFlowService.GetAllWorkFlows(c)
}

func createWorkFlow(c *gin.Context, workFlowService service.WorkFlowAbstractService) {
	var newWorkFlow dto.WorkFlowDTO
	if err := c.ShouldBindJSON(&newWorkFlow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	log.Println("Parsed workflow from request:", newWorkFlow)

	workFlowService.CreateWorkFlow(c, newWorkFlow)
}
