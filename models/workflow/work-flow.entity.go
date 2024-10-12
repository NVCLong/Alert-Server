package workflow

import (
	"time"

	"github.com/NVCLong/Alert-Server/common"
)

// WorkFlow represents a workflow entity
type WorkFlow struct {
	ID              uint
	UserFID         uint
	ConditionBatchs []ConditionBatch
	WorkFlowName    string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (c *WorkFlow) TableName() string {
	return "work_flows"
}

var _ common.Model = (*ConditionBatch)(nil)
