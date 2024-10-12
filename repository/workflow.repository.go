package repository

import (
	"context"

	"github.com/NVCLong/Alert-Server/models/workflow"
	"gorm.io/gorm"
)

type WorkFlowRepository struct {
	db *gorm.DB
}

func NewWorkFlowRepo(db *gorm.DB) *WorkFlowRepository {
	return &WorkFlowRepository{db: db}
}

func (tr *WorkFlowRepository) Create(c context.Context, workFlow workflow.WorkFlow) {

}
