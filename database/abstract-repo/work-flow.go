package abstractrepo

import (
	"github.com/NVCLong/Alert-Server/dto"
	"github.com/NVCLong/Alert-Server/models/workflow"
	"github.com/gin-gonic/gin"
)

const (
	WorkFlowTable  = "work_flows"
	ConditionBatch = "condition_batches"
)

type WorkFlowRepository interface {
	Create(c *gin.Context, workFlow *workflow.WorkFlow) error
	GetAll(c *gin.Context) ([]dto.GetAllWorkFlowRawResponse, error)
	GetById(id uint) (workflow.WorkFlow, error)
}
