package service

import (
	"fmt"
	"net/http"

	"github.com/NVCLong/Alert-Server/common"
	abstractrepo "github.com/NVCLong/Alert-Server/database/abstract-repo"
	"github.com/NVCLong/Alert-Server/dto"
	"github.com/NVCLong/Alert-Server/models/workflow"
	"github.com/gin-gonic/gin"
)

type WorkFlowService struct {
	repository abstractrepo.WorkFlowRepository
	logger     common.AbstractLogger
}

type WorkFlowAbstractService interface {
	GetAllWorkFlows(c *gin.Context)
	CreateWorkFlow(c *gin.Context, wf dto.WorkFlowDTO)
}

func NewWorkFlowService(repository abstractrepo.WorkFlowRepository) WorkFlowAbstractService {
	logger := common.NewTracingLogger("WorkFlowService")
	return &WorkFlowService{
		repository: repository,
		logger:     logger,
	}
}

func (s *WorkFlowService) GetAllWorkFlows(c *gin.Context) {
	s.logger.Debug("Starting query all work flows")
	result, err := s.repository.GetAll(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while get work flows"})
	}

	s.logger.Debug(fmt.Sprintf("Get total %d", len(result)))
	s.logger.Debug("Starting to map to get response")
	getResult := mapToWorkFlowResponse(result)

	s.logger.Debug("Mapping successfully")

	c.JSON(http.StatusOK, gin.H{"result": getResult})
}

func (s *WorkFlowService) CreateWorkFlow(c *gin.Context, wf dto.WorkFlowDTO) {
	var workFlow workflow.WorkFlow
	workFlow.UserFID = wf.UserID
	workFlow.WorkFlowName = wf.WorkFlowName
	s.repository.Create(c, &workFlow)
}

func mapToWorkFlowResponse(workFlows []dto.GetAllWorkFlowRawResponse) []dto.GetAllWorkFlowResponse {
	if len(workFlows) == 0 {
		return []dto.GetAllWorkFlowResponse{}
	}
	responses := []dto.GetAllWorkFlowResponse{}
	for _, workflow := range workFlows {
		response := dto.GetAllWorkFlowResponse{
			UserName:     workflow.UserName,
			WorkFlowId:   workflow.ID,
			WorkFlowName: workflow.WorkFlowName,
			LastUpdated:  common.FormatDate(workflow.UpdatedAt),
			CreatedAt:    common.FormatDate(workflow.CreatedAt),
		}

		responses = append(responses, response)
	}
	return responses
}
