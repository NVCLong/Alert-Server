package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NVCLong/Alert-Server/common"
	"github.com/NVCLong/Alert-Server/dto"
	"github.com/NVCLong/Alert-Server/modules/work-flow/repository"
	"github.com/NVCLong/Alert-Server/modules/work-flow/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewWorkFlowController(timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	// Define workflow route group
	workflowGroup := group.Group("/work-flow")

	// Initialize repository and services
	workFlowRepository := repository.NewWorkFlowRepository(db, "work_flows")

	// Register route handlers
	workflowGroup.GET("/all", func(ctx *gin.Context) {
		tracingID := common.GetTracingIDFromContext(ctx)
		logger := common.NewTracingLogger("WorkFlowController", tracingID)
		workFlowService := service.NewWorkFlowService(db, workFlowRepository, ctx)
		getAllWorkFlows(ctx, workFlowService, logger)
	})

	workflowGroup.POST("/create", func(ctx *gin.Context) {
		tracingID := common.GetTracingIDFromContext(ctx)
		logger := common.NewTracingLogger("WorkFlowController", tracingID)
		workFlowService := service.NewWorkFlowService(db, workFlowRepository, ctx)
		createWorkFlow(ctx, workFlowService, logger)
	})

	workflowGroup.POST("/import/:id", func(ctx *gin.Context) {
		tracingID := common.GetTracingIDFromContext(ctx)
		logger := common.NewTracingLogger("WorkFlowController", tracingID)
		workFlowService := service.NewWorkFlowService(db, workFlowRepository, ctx)
		importWorkFlow(ctx, workFlowService, logger)
	})

	workflowGroup.GET("/execute/:id", func(ctx *gin.Context) {
		tracingID := common.GetTracingIDFromContext(ctx)
		logger := common.NewTracingLogger("WorkFlowController", tracingID)
		workFlowService := service.NewWorkFlowService(db, workFlowRepository, ctx)
		executeWorkFlow(ctx, workFlowService, logger)
	})

	workflowGroup.DELETE("/:id", func(ctx *gin.Context) {
		tracingID := common.GetTracingIDFromContext(ctx)
		logger := common.NewTracingLogger("WorkFlowController", tracingID)
		workFlowService := service.NewWorkFlowService(db, workFlowRepository, ctx)
		deleteWorkFlow(ctx, workFlowService, logger)
	})

	workflowGroup.GET("/notifySync", func(ctx *gin.Context) {
		tracingID := common.GetTracingIDFromContext(ctx)
		logger := common.NewTracingLogger("WorkFlowController", tracingID)
		workFlowService := service.NewWorkFlowService(db, workFlowRepository, ctx)
		notifySync(ctx, workFlowService, logger)
	})
}

func notifySync(c *gin.Context, workFlowService service.WorkFlowAbstractService, logger common.AbstractLogger) {
	logger.Log("Start check and notify Sync Result")
	workFlowService.NotifySyncResult(c)
}

func getAllWorkFlows(c *gin.Context, workFlowService service.WorkFlowAbstractService, logger common.AbstractLogger) {
	logger.Log("Recieve Request get all workflows ")
	workFlowService.GetAllWorkFlows(c)
}

func deleteWorkFlow(c *gin.Context, workFlowService service.WorkFlowAbstractService, logger common.AbstractLogger) {
	logger.Log("Delete work flow")
	workFlowId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow id", "details": err.Error()})
		return
	}
	workFlowService.DeleteWorkFlow(c, uint(workFlowId))

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

func createWorkFlow(c *gin.Context, workFlowService service.WorkFlowAbstractService, logger common.AbstractLogger) {
	logger.Debug("Start to create workflow")
	var newWorkFlow dto.WorkFlowDTO
	if err := c.ShouldBindJSON(&newWorkFlow); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	log.Println("Parsed workflow from request:", newWorkFlow)

	workFlowService.CreateWorkFlow(c, newWorkFlow)
}

func importWorkFlow(c *gin.Context, workFlowService service.WorkFlowAbstractService, logger common.AbstractLogger) {
	logger.Debug("Start to import work flow")
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
