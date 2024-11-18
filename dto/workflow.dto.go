package dto

import (
	"github.com/NVCLong/Alert-Server/models/workflow"
	"time"
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

type ImoportWorkFlowRequest struct {
	ListCondition []WorkFlowConditionRequest
}

type WorkFlowConditionRequest struct {
	Condition string
	Action    string
}

type SyncEventDTO struct {
	SyncID         int64     `json:"sync_id" db:"sync_id"`
	SyncEvent      string    `json:"sync_event" db:"sync_event"`
	SyncStartTime  time.Time `json:"sync_start_time" db:"sync_start_time"`
	SyncFinishTime time.Time `json:"sync_finish_time" db:"sync_finish_time"`
	SyncStatus     bool      `json:"sync_status" db:"sync_status"`
	SyncFailReason string    `json:"sync_fail_reason,omitempty" db:"sync_fail_reason"`
	UserID         int64     `json:"user_id" db:"user_id"`
	Username       string    `json:"user_name"`
	UserEmail      string    `json:"user_email"`
	IsNew          bool      `json:"is_new" db:"is_new"`
}
