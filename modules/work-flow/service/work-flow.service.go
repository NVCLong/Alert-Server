package service

import (
	"net/http"

	"github.com/NVCLong/Alert-Server/common"
	abstractrepo "github.com/NVCLong/Alert-Server/database/abstract-repo"
	"github.com/NVCLong/Alert-Server/dto"
	"github.com/NVCLong/Alert-Server/models/workflow"
	"github.com/gin-gonic/gin"
)

type WorkFlowService struct {
	repository abstractrepo.WorkFlowRepository
}

type WorkFlowAbstractService interface {
	GetAllWorkFlows(c *gin.Context)
	CreateWorkFlow(c *gin.Context, wf dto.WorkFlowDTO)
}

func NewWorkFlowService(repository abstractrepo.WorkFlowRepository) WorkFlowAbstractService {
	return &WorkFlowService{
		repository: repository,
	}
}

func (s *WorkFlowService) GetAllWorkFlows(c *gin.Context) {
	result, err := s.repository.GetAll(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error while get work flows"})
	}
	getResult := mapToWorkFlowResponse(result)

	c.JSON(http.StatusOK, gin.H{"result": getResult})
}

func (s *WorkFlowService) CreateWorkFlow(c *gin.Context, wf dto.WorkFlowDTO) {
	var workFlow workflow.WorkFlow
	workFlow.UserFID = wf.UserID
	workFlow.WorkFlowName = wf.WorkFlowName
	s.repository.Create(c, &workFlow)
}

func mapToWorkFlowResponse(workFlows []dto.GetAllWorkFlowRawResponse) []dto.GetAllWorkFlowResponse {
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
