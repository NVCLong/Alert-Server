package dto

type ConditionBatchDTO struct {
	WorkFlowID  uint   `json:"workflow_id"`
	Condition   string `json:"condition"`
	BindingAttr string `json:"binding_attr"`
	Action      string `json:"action"`
	Type        string `json:"type"`
}

type SaveResponse struct {
	Status     bool
	Message    string
	WorkflowId uint
}

type JobCreation struct {
	WorkflowID string
	UserID     string
	Username   string
	UserEmail  string
}

type ParseResult struct {
	Operator        string
	Parameter       string
	Variable        string
	FunctionHandler string
}

type Deadline struct {
	Username    string
	Useremail   string
	Deadline    string
	Priority    string
	Description string
	Course      string
	Lecturer    string
}

type NotificationDeadlineRequest struct {
	Results []Deadline
	Action  string
	Type    string
}
