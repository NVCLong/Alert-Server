package conditionbatch

import (
	"github.com/NVCLong/Alert-Server/common"
	"github.com/NVCLong/Alert-Server/dto"
	"github.com/NVCLong/Alert-Server/models/workflow"
	"gorm.io/gorm"
	"time"
)

type AbstractService interface {
	SaveCondition(dto dto.ConditionBatchDTO) dto.SaveResponse
	GetCondition()
	TriggerCondition()
}
type ConditionBatchService struct {
	db     *gorm.DB
	logger common.AbstractLogger
}

func (c ConditionBatchService) SaveCondition(conDto dto.ConditionBatchDTO) dto.SaveResponse {
	conditionBatch := workflow.ConditionBatch{
		WorkFlowID:      conDto.WorkFlowID,
		Condition:       conDto.Condition,
		BindingAttr:     conDto.BindingAttr,
		SolvingFunction: "",
		Action:          conDto.Action,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := c.db.Create(&conditionBatch).Error; err != nil {
		return dto.SaveResponse{
			Status:     false,
			Message:    "Failed to save condition batch: " + err.Error(),
			WorkflowId: conDto.WorkFlowID,
		}
	}

	return dto.SaveResponse{
		Status:     true,
		Message:    "Condition Batch saved successfully",
		WorkflowId: conDto.WorkFlowID,
	}
}

func (c ConditionBatchService) GetCondition() {
	//TODO implement me
	panic("implement me")
}

func (c ConditionBatchService) TriggerCondition() {
	//TODO implement me
	panic("implement me")
}

func NewBatchService(db *gorm.DB, logger common.AbstractLogger) AbstractService {
	return &ConditionBatchService{
		db:     db,
		logger: logger,
	}
}
