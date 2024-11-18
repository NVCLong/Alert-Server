package service

import (
	"container/list"
	"encoding/json"
	"fmt"
	env "github.com/NVCLong/Alert-Server/bootstrap"
	conditionbatch "github.com/NVCLong/Alert-Server/modules/condition-batch"
	"github.com/NVCLong/Alert-Server/redis"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"sync"
	"time"

	"github.com/NVCLong/Alert-Server/common"
	abstractrepo "github.com/NVCLong/Alert-Server/database/abstract-repo"
	"github.com/NVCLong/Alert-Server/dto"
	"github.com/NVCLong/Alert-Server/models/workflow"
	"github.com/gin-gonic/gin"
)

type WorkFlowService struct {
	repository            abstractrepo.WorkFlowRepository
	userRepository        abstractrepo.UserAbstractRepository
	conditionBatchService conditionbatch.AbstractService
	logger                common.AbstractLogger
	dataSource            *gorm.DB
}

func (s *WorkFlowService) NotifySyncResult(c *gin.Context) {
	s.logger.Debug("Start querying sync event")

	var results []dto.SyncEventDTO

	// Define the SQL query for retrieving sync events and joined user data
	query := `SELECT 
                s.sync_id, 
                s.sync_event, 
                s.sync_start_time, 
                s.sync_finish_time, 
                s.sync_status, 
                s.sync_fail_reason, 
                s.user_id, 
                u.name AS username, 
                u.email AS user_email, 
                s.is_new
            FROM 
                sync_events s
            JOIN 
                student_users u ON s.user_id = u.id
            WHERE s.is_new = true;`

	if err := s.dataSource.Raw(query).Scan(&results).Error; err != nil {
		s.logger.Error("Failed to query sync events")
		c.JSON(500, gin.H{
			"message": "Failed to query sync events",
			"error":   err.Error(),
		})
		return
	}

	var successResults []dto.SyncEventDTO
	var failedResults []dto.SyncEventDTO

	for _, result := range results {
		if result.SyncStatus {
			successResults = append(successResults, result)
		} else {
			failedResults = append(failedResults, result)
		}
	}

	generateSyncNotification(results[0].UserEmail, results[0].Username, successResults, failedResults)

}

func generateSyncNotification(email string, name string, successResults []dto.SyncEventDTO, failedResults []dto.SyncEventDTO) {
	from := env.GetEnv(env.EnvEmail)
	pass := env.GetEnv(env.EnvPass)
	to := email
	fmt.Printf("From : %s To : %s\n", from, to)

	var sendContent string

	sendContent += "Send to " + name + "\n"
	sendContent += fmt.Sprintf("Total have %d sync events \n", len(successResults)+len(failedResults))
	sendContent += fmt.Sprintf("In this, %d events are success and %d events are fail \n", len(successResults), len(failedResults))
	sendContent += fmt.Sprintf("There is a list of fail events with reason: \n")

	for _, event := range failedResults {
		sendContent += fmt.Sprintf("Sync event : %s at %s with fail reason %s \n", event.SyncEvent, strings.Split(event.SyncStartTime.String(), " ")[0], event.SyncFailReason)
	}

	sendContent += "Daily sync report for admin"
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(sendContent))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
	fmt.Println("Success Send Email")
}

type WorkFlowAbstractService interface {
	GetAllWorkFlows(c *gin.Context)
	CreateWorkFlow(c *gin.Context, wf dto.WorkFlowDTO)
	ImportWorkFlow(c *gin.Context, workFlowId uint, listCondition []dto.WorkFlowConditionRequest) dto.SaveResponse
	ParseCondition(conditionObj dto.WorkFlowConditionRequest, logger common.AbstractLogger, workFlowId uint) (dto.SaveResponse, error)
	ExecuteWorkFlow(c *gin.Context, workFlowId uint, userIds []string)
	DeleteWorkFlow(c *gin.Context, workflowId uint)
	NotifySyncResult(c *gin.Context)
}

func NewWorkFlowService(db *gorm.DB, repository abstractrepo.WorkFlowRepository, ctx *gin.Context) WorkFlowAbstractService {
	tracingID := common.GetTracingIDFromContext(ctx)
	logger := common.NewTracingLogger("WorkFlowController", tracingID)
	batchService := conditionbatch.NewBatchService(db, ctx)
	userRepository := abstractrepo.NewUserRepository(db)
	return &WorkFlowService{
		repository:            repository,
		logger:                logger,
		conditionBatchService: batchService,
		userRepository:        userRepository,
		dataSource:            db,
	}
}

func (s *WorkFlowService) DeleteWorkFlow(c *gin.Context, workflowId uint) {
	s.logger.Debug("Starting to delete WorkFlow")
	if err := s.dataSource.Delete(&workflow.WorkFlow{}, workflowId).Error; err != nil {
		s.logger.Error(err.Error())
		c.JSON(500, gin.H{
			"error": "Failed to delete WorkFlow",
		})
		return
	}
	s.conditionBatchService.DeleteConditionByWorkFlowId(workflowId)
	s.logger.Debug(fmt.Sprintf("Successfully deleted WorkFlow with ID: ", workflowId))
	common.SetTraceIDHeader(c, s.logger.GetTraceId())
	c.JSON(200, gin.H{
		"message": "WorkFlow deleted successfully",
	})
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
	common.SetTraceIDHeader(c, s.logger.GetTraceId())
	c.JSON(http.StatusOK, gin.H{"result": getResult})
}

func (s *WorkFlowService) ExecuteWorkFlow(c *gin.Context, workFlowId uint, userIds []string) {
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 7)
	s.logger.Debug(fmt.Sprintf("Trigger create notification from %s to %s", startDate, endDate))
	results := s.userRepository.FindUserHaveDeadline(startDate, endDate)
	if len(results) == 0 {
		s.logger.Debug("Do not have any user have deadline in next week")
		c.JSON(http.StatusOK, gin.H{"response": "Do not have user have deadline in this week"})
		return
	}
	redisClient := redis.NewRedisConnection()
	defer redisClient.Close()

	var wg sync.WaitGroup

	// start worker pool
	redis.StartWorkerPool(redisClient, &wg, c, len(results), s.dataSource)
	for _, result := range results {
		job := fmt.Sprintf(`{ "WorkflowID": "%d", "UserID": "%d", "Username": "%s", "UserEmail": "%s" }`, workFlowId, result.Id, result.Name, result.Email)
		redis.PushJobToQueue(c, job, s.logger)
	}
	// Wait for all workers to finish
	wg.Wait()
	common.SetTraceIDHeader(c, s.logger.GetTraceId())
	c.JSON(http.StatusOK, gin.H{"message": "Process success"})
}

func (s *WorkFlowService) CreateWorkFlow(c *gin.Context, wf dto.WorkFlowDTO) {
	var workFlow workflow.WorkFlow
	workFlow.UserFID = wf.UserID
	workFlow.WorkFlowName = wf.WorkFlowName
	s.repository.Create(c, &workFlow)
}

func (s *WorkFlowService) ImportWorkFlow(c *gin.Context, workFlowId uint, listCondition []dto.WorkFlowConditionRequest) dto.SaveResponse {
	workFlow, err := s.repository.GetById(workFlowId)
	if err != nil {
		return dto.SaveResponse{
			Message:    "Workflow is not found",
			WorkflowId: workFlowId,
		}
	}
	s.conditionBatchService.DeleteConditionByWorkFlowId(workFlowId)
	s.logger.Debug(fmt.Sprintf("Found Workflow with id : %d", workFlow.ID))
	for _, condition := range listCondition {
		result, error := s.ParseCondition(condition, s.logger, workFlowId)
		if error != nil || !result.Status {
			return dto.SaveResponse{
				Message:    "Fail to import work flow",
				Status:     false,
				WorkflowId: workFlowId,
			}
		}
	}
	s.logger.Debug(fmt.Sprintf("Import condition to workflow %d", workFlowId))
	common.SetTraceIDHeader(c, s.logger.GetTraceId())
	return dto.SaveResponse{
		Message:    "Import success",
		Status:     true,
		WorkflowId: workFlowId,
	}
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

// condition should be:
// AT LEAST OF ONE DEADLINE IN [.......] with ..... is parameter
// DEADLINE TIME  < Date AND DEADLINE >DATE
// TOTAL DEADLINE IN WEEK
// TOTAL DEADLINE FINISH IN WEEK
// DEADLINE NEED TO DONE IN ? DAYS AFTER
// CONDITION A AND CONDITION B
func (s *WorkFlowService) ParseCondition(conditionObj dto.WorkFlowConditionRequest, logger common.AbstractLogger, workFlowId uint) (dto.SaveResponse, error) {
	// Create a map to store results for each sub-condition\
	condition := conditionObj.Condition
	var resultsArray []*dto.ParseResult
	var conditionString string
	conditionKeywords := list.New()
	isSingleCondition := getAndCheckConditionKeyword(condition, conditionKeywords, logger)

	// Loop through each condition keyword to find matches condition have condition keyword
	if !isSingleCondition {
		keyword := conditionKeywords.Front().Value.(string)
		if strings.Contains(condition, keyword) {
			conditionKeywords.PushBack(keyword)
			subConditions := strings.Split(condition, keyword)
			// Parse each sub-condition and accumulate the results in resultMap
			parsedResults := parseSubCondition(subConditions, *conditionKeywords, logger)

			for _, result := range parsedResults {
				resultsArray = append(resultsArray, result)
				conditionKeywords.PushBack(result.FunctionHandler)
			}
		}
		conditionOperatorStrings := strings.Join(common.CONDITION_KEYWORD, ",")
		conditionOperator := list.New()
		var resultArray []string

		if conditionKeywords.Len() != 0 {
			for e := conditionKeywords.Front(); e != nil; {
				// If the current value is a condition keyword
				if val, ok := e.Value.(string); ok && strings.Contains(conditionOperatorStrings, val) {
					// Store the condition operator
					conditionOperator.PushBack(val)
				} else {
					if len(resultArray) == 0 {
						resultArray = append(resultArray, val)
					} else {
						if conditionOperator.Len() > 0 {
							operatorElement := conditionOperator.Remove(conditionOperator.Back()).(string)
							resultArray = append(resultArray, operatorElement, val)
						} else {
							resultArray = append(resultArray, val)
						}
					}
				}
				e = e.Next()
			}
			conditionString = strings.Join(resultArray, " ")
		}
	}

	// case condition do not contains condition keyword
	if isSingleCondition {
		parsedResults := parseSingleCondition(condition, logger)
		for _, result := range parsedResults {
			resultsArray = append(resultsArray, result)
			conditionString = result.FunctionHandler
		}
	}

	logger.Debug(conditionString)
	logger.Debug("Convert the array to JSON")
	// Convert the array to JSON
	jsonData, err := json.Marshal(resultsArray)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to convert results to JSON: %s", err))
		response := dto.SaveResponse{
			Message:    "Import Workflow fail as fail to convert results to JSON",
			Status:     false,
			WorkflowId: workFlowId,
		}
		return response, err
	}

	var action Action

	// Parse the JSON string
	correctedJSON := strings.ReplaceAll(conditionObj.Action, "'", "\"")
	err = json.Unmarshal([]byte(correctedJSON), &action)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		response := dto.SaveResponse{
			Message:    "Import Workflow fail as fail to convert results to JSON",
			Status:     false,
			WorkflowId: workFlowId,
		}
		return response, err
	}

	// binding attribute
	jsonString := string(jsonData)

	//condition

	logger.Debug(jsonString)

	// Return the map containing all the parsed results
	saveConditionDto := dto.ConditionBatchDTO{
		WorkFlowID:  workFlowId,
		Condition:   conditionString,
		BindingAttr: jsonString,
		Action:      conditionObj.Action,
		Type:        action.Type,
	}
	result := s.conditionBatchService.SaveCondition(saveConditionDto)
	return result, nil
}

func getAndCheckConditionKeyword(condition string, conditionKeyWords *list.List, logger common.AbstractLogger) bool {
	for _, keyword := range common.CONDITION_KEYWORD {
		if strings.Contains(condition, keyword) {
			logger.Debug("Have condition keyword")
			conditionKeyWords.PushBack(keyword)
			return false
		}
	}
	return true
}

func parseSingleCondition(condition string, logger common.AbstractLogger) map[string]*dto.ParseResult {
	operatorMap := map[OperatorType][]string{
		NUMERICS: common.NUMERIC_OPERATOR,
		EQUALS:   common.EQUAL_OPERATOR,
		CONTAINS: common.CONTAIN_OPERATOR,
	}

	resultMap := make(map[string]*dto.ParseResult)

	foundOperator := false
	for operatorType, operators := range operatorMap {
		for _, operator := range operators {
			if strings.Contains(condition, operator) {
				logger.Debug(fmt.Sprintf("Found %s operator: %s in subCondition: %s", operatorType, operator, condition))
				foundOperator = true

				// Parse the condition and store the result in the map
				parseResult := parseConditionToFunction(condition, operator, operatorType, logger)
				resultMap[condition] = parseResult

				logger.Debug(fmt.Sprintf("Stored result for subCondition: %s", condition))
				break
			}
		}
		if foundOperator {
			break
		}
	}

	return resultMap
}

func parseSubCondition(subConditions []string, conditionKeywords list.List, logger common.AbstractLogger) map[string]*dto.ParseResult {

	operatorMap := map[OperatorType][]string{
		NUMERICS: common.NUMERIC_OPERATOR,
		EQUALS:   common.EQUAL_OPERATOR,
		CONTAINS: common.CONTAIN_OPERATOR,
	}

	resultMap := make(map[string]*dto.ParseResult)

	for _, subCondition := range subConditions {
		foundOperator := false
		for operatorType, operators := range operatorMap {
			for _, operator := range operators {
				if strings.Contains(subCondition, operator) {
					logger.Debug(fmt.Sprintf("Found %s operator: %s in subCondition: %s", operatorType, operator, subCondition))
					foundOperator = true

					// Parse the condition and store the result in the map
					parseResult := parseConditionToFunction(subCondition, operator, operatorType, logger)
					resultMap[subCondition] = parseResult

					logger.Debug(fmt.Sprintf("Stored result for subCondition: %s", subCondition))
					break
				}
			}
			if foundOperator {
				break
			}
		}
	}

	return resultMap

}

func parseConditionToFunction(contextSring string, operator string, opType OperatorType, logger common.AbstractLogger) *dto.ParseResult {
	logger.Debug("Start to parse condition to function context")
	params := strings.Split(contextSring, operator)[1]
	compareString := strings.Split(contextSring, operator)[0]
	logger.Debug(fmt.Sprintf("Found Params: %s", params))
	logger.Debug(fmt.Sprintf("Found compare string: %s", compareString))
	functionHanlderName := getFunctionHandlerName(compareString, operator, params)
	logger.Debug(fmt.Sprintf("variable: %s, Functionhanler: %s", functionHanlderName.variable, functionHanlderName.handler))

	return &dto.ParseResult{
		Operator:        operator,
		Parameter:       params,
		Variable:        functionHanlderName.variable,
		FunctionHandler: functionHanlderName.handler,
	}
}

func getFunctionHandlerName(compareString string, opeator string, params string) *GetHandlerResposne {
	if strings.Contains(compareString, common.AT_LEAST_PATTERN) {
		variable := strings.ReplaceAll(compareString, common.AT_LEAST_PATTERN, "")
		return &GetHandlerResposne{
			variable: variable,
			handler:  string(common.AT_LEAST_HANDLER),
		}
	}

	if strings.Contains(compareString, common.TOTAL_PATERN) {
		variable := strings.ReplaceAll(compareString, common.TOTAL_PATERN, "")
		return &GetHandlerResposne{
			variable: variable,
			handler:  string(common.TOTAL_HANDLER),
		}
	}

	if strings.Contains(compareString, common.HAVE_MORE_PATERN) {
		variable := strings.ReplaceAll(compareString, common.HAVE_MORE_PATERN, "")
		return &GetHandlerResposne{
			variable: variable,
			handler:  string(common.HAVE_MORE_HANDLER),
		}
	}

	if strings.Contains(compareString, common.HAVE) {
		variable := strings.ReplaceAll(compareString, common.HAVE, "")
		return &GetHandlerResposne{
			variable: variable,
			handler:  string(common.HAVE_HANDLER),
		}
	}

	return nil
}

type OperatorType string

const (
	NUMERICS OperatorType = "numerics"
	EQUALS   OperatorType = "equals"
	CONTAINS OperatorType = "contains"
)

type GetHandlerResposne struct {
	variable string
	handler  string
}

type Action struct {
	Message string
	Type    string
}
