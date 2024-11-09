package conditionbatch

import (
	"encoding/json"
	"fmt"
	env "github.com/NVCLong/Alert-Server/bootstrap"
	"github.com/NVCLong/Alert-Server/common"
	abstractrepo "github.com/NVCLong/Alert-Server/database/abstract-repo"
	"github.com/NVCLong/Alert-Server/dto"
	"github.com/NVCLong/Alert-Server/models/workflow"
	"gorm.io/gorm"
	"log"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

type AbstractService interface {
	SaveCondition(dto dto.ConditionBatchDTO) dto.SaveResponse
	GetCondition(workflowID uint) dto.ConditionBatchDTO
	TriggerCondition(jobCreation dto.JobCreation)
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

func (c ConditionBatchService) GetCondition(workFlowID uint) dto.ConditionBatchDTO {
	c.logger.Log(fmt.Sprintf("Start to find condition batch with workFlowId : %s", strconv.Itoa(int(workFlowID))))
	var conditionBatch dto.ConditionBatchDTO
	err := c.db.Table(abstractrepo.ConditionBatch).Where("work_flow_id = ?", workFlowID).First(&conditionBatch).Error
	if err != nil {
		c.logger.Debug("Can not get condition batch with work flow id")
	}
	return conditionBatch
}

func (c ConditionBatchService) TriggerCondition(jobCreation dto.JobCreation) {
	c.logger.Log(fmt.Sprintf("Start to execute workflow %s with userid %s", jobCreation.WorkflowID, jobCreation.UserID))
	workFlowID, err := strconv.ParseUint(jobCreation.WorkflowID, 10, 64)
	userID, err := strconv.ParseUint(jobCreation.UserID, 10, 64)
	if err != nil {
		c.logger.Log(fmt.Sprintf("Failed to parse workFlowID: %v", err))
		return
	}
	conditionBatch := c.GetCondition(uint(workFlowID))
	c.logger.Debug(fmt.Sprintf("Find Condition Batch %s", conditionBatch.BindingAttr))
	var bindingConditions []dto.ParseResult
	json.Unmarshal([]byte(conditionBatch.BindingAttr), &bindingConditions)

	for _, result := range bindingConditions {
		c.ExecuteFunctionHandler(result, uint(userID), conditionBatch.Action)
	}

}

func (s ConditionBatchService) ExecuteFunctionHandler(result dto.ParseResult, userId uint, action string) {
	switch result.FunctionHandler {
	case string(common.AT_LEAST_HANDLER):
		handleAtLeastCase(result.Operator, result.Parameter, result.Variable, userId, s.db, action)
	default:
		return
	}
}

func handleAtLeastCase(operator string, parameter string, variable string, userId uint, db *gorm.DB, action string) {
	var result []dto.Deadline
	if includes(common.DATES, parameter).isInDate {
		query := fmt.Sprintf(`
		SELECT
		  deadline.deadline as Deadline, deadline.priority as Priority, deadline.description as Description, 
		  courses.course_name as Course, course_value.lecture as Lecturer, student_users.name as Username, 
      	  student_users.email as Useremail
		FROM deadline
		  JOIN course_value ON deadline."courseValueId" = course_value.course_value_id
		  JOIN courses ON courses.course_id = course_value."coursesId"
          JOIN student_users on student_users.id= deadline."userId"
		WHERE
		  deadline."userId" = %d
		  and deadline."is_Active" = true
		  AND deadline.deadline::date %s (DATE '%s');`, int(userId), operator, strings.Split(includes(common.DATES, parameter).date.String(), " ")[0])
		db.Raw(query).Scan(&result)
		if len(result) == 0 {
			return
		}
		notiReq := dto.NotificationDeadlineRequest{
			Results: result,
			Action:  action,
		}
		generateDeadlineNotification(notiReq, variable)
	}
}

func generateDeadlineNotification(notiReq dto.NotificationDeadlineRequest, variable string) {
	from := env.GetEnv(env.EnvEmail)
	pass := env.GetEnv(env.EnvPass)
	to := notiReq.Results[0].Useremail
	fmt.Printf("From : %s To : %s", from, to)
	var sendContent string
	sendContent = "Subject: Weekly Notification\n\n"
	sendContent += fmt.Sprintf("Send to %s \n", notiReq.Results[0].Username)
	for _, result := range notiReq.Results {
		sendContent += fmt.Sprintf("You have one %s in course %s with deadline %s with %s priority \n", variable, result.Course, result.Deadline, result.Priority)
	}
	fmt.Println(sendContent)
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(sendContent))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
	fmt.Println("Success Send Email")
}

func includes(slice []string, item string) GetDayFromDate {
	item = strings.TrimSpace(strings.ToUpper(item))
	today := time.Now()
	for _, v := range slice {
		if v == item {
			daysOffset := common.DayOffsets[item]
			todayOffset := common.DayOffsets[today.String()[:3]]
			daysUntil := (daysOffset - todayOffset + 7) % 7
			if daysUntil == 0 {
				daysUntil = 7
			}
			nextDay := today.AddDate(0, 0, daysUntil)
			return GetDayFromDate{
				isInDate: true,
				date:     nextDay,
			}
		}
	}
	return GetDayFromDate{
		isInDate: false,
	}
}
func NewBatchService(db *gorm.DB) AbstractService {
	logger := common.NewTracingLogger("BatchService")
	return &ConditionBatchService{
		db:     db,
		logger: logger,
	}
}

type GetDayFromDate struct {
	isInDate bool
	date     time.Time
}
