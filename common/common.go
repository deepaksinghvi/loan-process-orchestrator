package common

const (
	LOTaskQueueName       = "lo-task"
	SignalName            = "loan-signal"
	LoanWorkflowQueryType = "loan-workflow-state-query"
)

type State string
type QueryResult struct {
	State   State
	Content string
}

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

type WorkflowObject struct {
	WorkflowID string `json:"workflow_id"`
	RunID      string `json:"run_id"`
}

type ApprovalWorkflowRequest struct {
	WorkflowID string `json:"workflow_id"`
	RunID      string `json:"run_id"`
	Approved   bool   `json:"approved"`
}
