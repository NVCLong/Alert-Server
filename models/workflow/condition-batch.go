package workflow

import (
	"time"

	"github.com/NVCLong/Alert-Server/common"
)

// ConditionBatch represents a condition batch entity
type ConditionBatch struct {
	ID              uint   `gorm:"primaryKey"` // Primary key
	WorkFlowID      uint   // Foreign key referencing WorkFlow
	Condition       string // Fixed typo from Condinmtion to Condition
	SolvingFunction string
	BindingAttr     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (c *ConditionBatch) TableName() string {
	return "condition_batches"
}

var _ common.Model = (*ConditionBatch)(nil)
