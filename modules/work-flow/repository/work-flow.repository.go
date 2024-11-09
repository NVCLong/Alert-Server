package repository

import (
	abstractrepo "github.com/NVCLong/Alert-Server/database/abstract-repo"
	"github.com/NVCLong/Alert-Server/dto"
	"github.com/NVCLong/Alert-Server/models/workflow"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkFlowRepository struct {
	database *gorm.DB
	table    string
}

func (r *WorkFlowRepository) GetById(id uint) (workflow.WorkFlow, error) {
	var workFlowRes workflow.WorkFlow
	result := r.database.Table(r.table).Where("id = ?", id).First(&workFlowRes)

	if result.Error != nil {
		return workFlowRes, result.Error
	}

	return workFlowRes, nil
}

func NewWorkFlowRepository(db *gorm.DB, table string) abstractrepo.WorkFlowRepository {
	return &WorkFlowRepository{
		database: db,
		table:    table,
	}
}

func (r *WorkFlowRepository) Create(c *gin.Context, workFlow *workflow.WorkFlow) error {
	result := r.database.Table(r.table).Create(&workFlow)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *WorkFlowRepository) GetAll(c *gin.Context) ([]dto.GetAllWorkFlowRawResponse, error) {
	var workflows []dto.GetAllWorkFlowRawResponse

	result := r.database.WithContext(c).Table("work_flows").Select("work_flows.*, student_users.name AS username").Joins("JOIN student_users ON student_users.id = work_flows.user_f_id").Find(&workflows)
	if result.Error != nil {
		return nil, result.Error
	}

	return workflows, nil
}
