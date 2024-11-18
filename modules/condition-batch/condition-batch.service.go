package conditionbatch

import (
	"encoding/json"
	"fmt"
	env "github.com/NVCLong/Alert-Server/bootstrap"
	"github.com/NVCLong/Alert-Server/common"
	abstractrepo "github.com/NVCLong/Alert-Server/database/abstract-repo"
	"github.com/NVCLong/Alert-Server/dto"
	"github.com/NVCLong/Alert-Server/models/workflow"
	"github.com/gin-gonic/gin"
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
	DeleteConditionByWorkFlowId(workflowId uint)
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
		Type:            conDto.Type,
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
	case string(common.TOTAL_HANDLER):
		handleTotalCase(result.Operator, result.Parameter, result.Variable, userId, s.db, action)
	default:
		return
	}
}

func handleAtLeastCase(operator string, parameter string, variable string, userId uint, db *gorm.DB, action string) {
	var result []dto.Deadline
	if includes(common.DATES, parameter).isInDate {
		dates := includes(common.DATES, parameter).dates
		dateList := make([]string, len(dates))
		for i, date := range dates {
			dateList[i] = fmt.Sprintf("DATE '%s'", strings.Split(date.String(), " ")[0])
		}
		dateString := strings.Join(dateList, ", ")
		fmt.Println(dateString)
		query := fmt.Sprintf(`
		SELECT
		  deadline.deadline AS Deadline, 
		  deadline.priority AS Priority, 
		  deadline.description AS Description, 
		  courses.course_name AS Course, 
		  course_value.lecture AS Lecturer, 
		  student_users.name AS Username, 
		  student_users.email AS Useremail
		FROM deadline
		  JOIN course_value ON deadline."courseValueId" = course_value.course_value_id
		  JOIN courses ON courses.course_id = course_value."coursesId"
          JOIN student_users ON student_users.id = deadline."userId"
		WHERE
		  deadline."userId" = %d
		  AND deadline."is_Active" = TRUE
		  AND deadline.deadline::date %s (%s);`, int(userId), operator, dateString)
		db.Raw(query).Scan(&result)
		if len(result) == 0 {
			fmt.Println("Do not have any deadline in this day")
			return
		}
		notiReq := dto.NotificationDeadlineRequest{
			Results: result,
			Action:  action,
		}
		generateDeadlineNotification(notiReq, variable)
	}
}
func handleTotalCase(operator string, parameter string, variable string, userId uint, db *gorm.DB, action string) {
	minDate := strings.Split(time.Now().String(), " ")[0]
	maxDate := strings.Split(getDateRange(parameter).String(), " ")[0]
	var result []dto.Deadline
	query := fmt.Sprintf(`
		SELECT
		  deadline.deadline AS Deadline, 
		  deadline.priority AS Priority, 
		  deadline.description AS Description, 
		  courses.course_name AS Course, 
		  course_value.lecture AS Lecturer, 
		  student_users.name AS Username, 
		  student_users.email AS Useremail
		FROM deadline
		  JOIN course_value ON deadline."courseValueId" = course_value.course_value_id
		  JOIN courses ON courses.course_id = course_value."coursesId"
          JOIN student_users ON student_users.id = deadline."userId"
		WHERE
		  deadline."userId" = %d
		  AND deadline."is_Active" = TRUE
		  AND deadline.deadline::date >= (DATE '%s')
		  AND deadline.deadline::date < (DATE '%s');`, int(userId), minDate, maxDate)

	db.Raw(query).Scan(&result)
	if len(result) == 0 {
		return
	}

	notiReq := dto.NotificationDeadlineRequest{
		Results: result,
		Action:  action,
		Type:    parameter,
	}

	generateDeadlineNotification(notiReq, variable)
}

func generateDeadlineNotification(notiReq dto.NotificationDeadlineRequest, variable string) {
	from := env.GetEnv(env.EnvEmail)
	pass := env.GetEnv(env.EnvPass)
	to := notiReq.Results[0].Useremail
	fmt.Printf("From : %s To : %s\n", from, to)

	sendContent := generateEmailContent(notiReq, variable)

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

func generateEmailContent(notiReq dto.NotificationDeadlineRequest, variable string) string {
	var sendContent string
	var action Action
	actionString := strings.Replace(notiReq.Action, "'", "\"", -1)
	err := json.Unmarshal([]byte(actionString), &action)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}
	switch notiReq.Type {
	case common.ALERT_TYPE_WEEKLY:
		sendContent = "Subject: Weekly Notification\n"
		sendContent += fmt.Sprintf("Send to %s\n", notiReq.Results[0].Username)
		sendContent += action.Message + "\n"
		if action.Type == strings.TrimSpace(common.AT_LEAST_PATTERN) {
			for _, result := range notiReq.Results {
				sendContent += fmt.Sprintf("You have one %s in course %s with deadline %s with %s priority\n",
					variable, result.Course, result.Deadline, result.Priority)
			}
		} else if action.Type == strings.TrimSpace(common.TOTAL_PATERN) {
			sendContent += fmt.Sprintf("You have to total %d %s in this week \n", len(notiReq.Results), variable)
			for _, result := range notiReq.Results {
				sendContent += fmt.Sprintf("You have one %s in course %s with deadline %s with %s priority\n",
					variable, result.Course, result.Deadline, result.Priority)
			}
		}
	case common.ALERT_TYPE_DAILY:
		sendContent = "Subject: Daily Notification\n"
		sendContent += fmt.Sprintf("Send to %s\n", notiReq.Results[0].Username)
		sendContent += action.Message + "\n"
		if action.Type == strings.TrimSpace(common.AT_LEAST_PATTERN) {
			for _, result := range notiReq.Results {
				sendContent += fmt.Sprintf("You have one %s in course %s with deadline %s with %s priority\n ",
					variable, result.Course, result.Deadline, result.Priority)
			}
		} else if action.Type == strings.TrimSpace(common.TOTAL_PATERN) {
			sendContent += fmt.Sprintf("You have to total %d %s today \n", len(notiReq.Results), variable)
			for _, result := range notiReq.Results {
				sendContent += fmt.Sprintf("You have one %s in course %s with deadline %s with %s priority\n",
					variable, result.Course, strings.Split(result.Deadline, " ")[0], result.Priority)
			}
		}
	case common.ALERT_TYPE_SCHEDULE:
		sendContent = "Subject: Scheduler Notification\n"
		sendContent += fmt.Sprintf("Send to %s\n", notiReq.Results[0].Username)
	default:
		sendContent = "Subject: General Notification\n"
		sendContent += fmt.Sprintf("Send to %s\n", notiReq.Results[0].Username)
	}

	return sendContent
}

func (s *ConditionBatchService) DeleteConditionByWorkFlowId(workflowId uint) {
	if err := s.db.Delete(&workflow.ConditionBatch{}).Where("work_flow_f_id= ?", workflowId).Error; err != nil {
		s.logger.Error("Fail to delete condition")
	}
}

func getDateRange(parameter string) time.Time {
	fmt.Println(parameter)
	dayOffsets := common.TotalOffset[strings.TrimSpace(parameter)]
	fmt.Println(dayOffsets)
	today := time.Now()
	maxDay := today.AddDate(0, 0, dayOffsets)
	fmt.Println(maxDay)
	return maxDay
}
func includes(slice []string, item string) GetDayFromDate {
	item = strings.TrimSpace(strings.ToUpper(item))
	days := strings.Split(item, ",")
	today := time.Now()
	var parameterDays []time.Time
	for _, day := range days {
		for _, v := range slice {
			if v == day {
				daysOffset := common.DayOffsets[day]
				todayOffset := common.DayOffsets[today.String()[:3]]
				daysUntil := (daysOffset - todayOffset + 7) % 7
				if daysUntil == 0 {
					daysUntil = 7
				}
				nextDay := today.AddDate(0, 0, daysUntil)
				parameterDays = append(parameterDays, nextDay)
			}
		}
	}
	if len(parameterDays) != 0 {
		return GetDayFromDate{
			isInDate: true,
			dates:    parameterDays,
		}
	}

	return GetDayFromDate{
		isInDate: false,
	}
}
func NewBatchService(db *gorm.DB, ctx *gin.Context) AbstractService {
	tracingID := common.GetTracingIDFromContext(ctx)
	logger := common.NewTracingLogger("BatchService", tracingID)
	return &ConditionBatchService{
		db:     db,
		logger: logger,
	}
}

type GetDayFromDate struct {
	isInDate bool
	dates    []time.Time
}

type Action struct {
	Message string
	Type    string
}
