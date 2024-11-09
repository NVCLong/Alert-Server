package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NVCLong/Alert-Server/common"
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
	logger := common.NewTracingLogger("WorkflowController")

	workFlowService := service.NewWorkFlowService(db, workFlowRepository)

	workflowGroup.GET("/all", func(ctx *gin.Context) {
		getAllWorkFlows(ctx, workFlowService, logger)
	})
	workflowGroup.POST("/create", func(ctx *gin.Context) {
		createWorkFlow(ctx, workFlowService)
	})
	workflowGroup.POST("/import/:id", func(ctx *gin.Context) {
		importWorkFlow(ctx, workFlowService)
	})
	workflowGroup.GET("/excute/:id", func(ctx *gin.Context) {
		executeWorkFlow(ctx, workFlowService, logger)
	})
}

func getAllWorkFlows(c *gin.Context, workFlowService service.WorkFlowAbstractService, logger common.AbstractLogger) {
	logger.Log("Recieve Request get all workflows ")
	workFlowService.GetAllWorkFlows(c)
}

func executeWorkFlow(c *gin.Context, workFlowService service.WorkFlowAbstractService, logger common.AbstractLogger) {
	logger.Log("Recieve Request execute work-flow")
	workFlowId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow id", "details": err.Error()})
		return
	}

	userIds := c.Query("userIds")
	userIdList := strings.Split(userIds, ",")
	workFlowService.ExecuteWorkFlow(c, uint(workFlowId), userIdList)

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

func importWorkFlow(c *gin.Context, workFlowService service.WorkFlowAbstractService) {
	var importRequest dto.ImoportWorkFlowRequest

	workFlowId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow id", "details": err.Error()})
		return
	}

	if err := c.ShouldBindBodyWithJSON(&importRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workFlowService.ImportWorkFlow(c, uint(workFlowId), importRequest.ListCondition))
}
