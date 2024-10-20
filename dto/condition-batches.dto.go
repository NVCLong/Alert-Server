package dto

type ConditionBatchDTO struct {
	WorkFlowID  uint   `json:"workflow_id"`
	Condition   string `json:"condition"`
	BindingAttr string `json:"binding_attr"`
	Action      string `json:"action"`
}

type SaveResponse struct {
	Status     bool
	Message    string
	WorkflowId uint
}
