package dto

import (
	"github.com/NVCLong/Alert-Server/models/workflow"
)

type WorkFlowDTO struct {
	UserID       uint
	WorkFlowName string
}

type GetAllWorkFlowRawResponse struct {
	workflow.WorkFlow
	UserName string `json:"username" gorm:"column:username"`
}

type GetAllWorkFlowResponse struct {
	UserName     string
	WorkFlowId   uint
	WorkFlowName string
	LastUpdated  string
	CreatedAt    string
}
