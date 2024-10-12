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

	result := r.database.WithContext(c).Table("work_flows").Select("work_flows.*, users.name AS username").Joins("JOIN users ON users.id = work_flows.user_f_id").Find(&workflows)
	if result.Error != nil {
		return nil, result.Error
	}

	return workflows, nil
}
