package controller

import (
	"github.com/deepaksinghvi/loan-process-orchestrator/common"
	"github.com/deepaksinghvi/loan-process-orchestrator/common/dto"
	"github.com/deepaksinghvi/loan-process-orchestrator/worker"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// CreateLoanApplication godoc
// @Summary Create New Loan APplication
// @ID create-loan-application
// @Tags loan-origination
// @Accept json
// @Produce json
// @Param loanApplicationRequest body dto.LoanApplicationInputStep true "Loan Application "
// @Success 200 {object} common.WorkflowObject
// @Failure 500 {object} common.HTTPError
// @Router /loan-application [post]
func CreateLoanApplication(c *gin.Context) {
	var loanApplicationInput dto.LoanApplicationInputStep
	// Call BindJSON to bind the received JSON to input.
	if err := c.ShouldBindJSON(&loanApplicationInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workflowID, runID, err := worker.StartLoanWorkflowExecution(loanApplicationInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"data": err.Error()})
	}
	log.Infof("Loan Application Initiated, WorkflowID: %s, RunID:%s", workflowID, runID)

	c.JSON(http.StatusOK, gin.H{"data": common.WorkflowObject{
		WorkflowID: workflowID,
		RunID:      runID,
	}})
}

// GetLoanApplication godoc
// @Summary GetLoan Application state
// @ID get-loan-application
// @Tags loan-origination
// @Accept json
// @Produce json
// @Param workflow_id path string true "Workflow ID"
// @Param run_id path string true "Run ID"
// @Success 200 {object} common.QueryResult
// @Failure 500 {object} common.HTTPError
// @Router /loan-application/{workflow_id}/{run_id} [get]
func GetLoanApplication(c *gin.Context) {
	workflowID := c.Param("workflow_id")
	runID := c.Param("run_id")
	result, err := worker.QueryWorkflow(workflowID, runID, common.LoanWorkflowQueryType)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"workflow state": result})
}

// LoanApplicationApproval godoc
// @Summary Loan Application Process
// @ID approval-process-loan-application
// @Tags loan-origination
// @Accept json
// @Produce json
// @Param approvalInput body common.ApprovalWorkflowRequest true "Loan Application Approval"
// @Success 200 {object} string
// @Failure 500 {object} common.HTTPError
// @Router /loan-application-approval/{workflow_id}/{run_id} [post]
func LoanApplicationApproval(c *gin.Context) {
	var approvalInput common.ApprovalWorkflowRequest
	// Call BindJSON to bind the received JSON to input.
	if err := c.ShouldBindJSON(&approvalInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	approvalSingalValue := "REJECTED"
	if approvalInput.Approved {
		approvalSingalValue = "APPROVED"
	}

	log.Infof("User Input Recevied for approval process as %t for workflowID: %s, RunID: %s",
		approvalInput.Approved, approvalInput.WorkflowID, approvalInput.RunID)
	err := worker.SendSignalToWorkflow(approvalInput.WorkflowID, approvalInput.RunID, common.SignalName, approvalSingalValue)
	if nil != err {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error during approval process!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "Approval Process Completed"})
}

/*
func queryLoanOriginationApplicationState(workflowID, runID string, args ...interface{}) (*common.QueryResult, error) {
	var result common.QueryResult
	workflowClient, err := common.LOCadenceHelper.Builder.BuildCadenceClient()
	resp, err := workflowClient.QueryWorkflowWithOptions(context.Background(),
		&client.QueryWorkflowWithOptionsRequest{
			WorkflowID:            workflowID,
			RunID:                 runID,
			QueryType:             common.QueryNameLoWorkflowState,
			QueryConsistencyLevel: shared.QueryConsistencyLevelStrong.Ptr(),
			Args:                  args,
		})
	if err != nil {
		common.LOCadenceHelper.Logger.Error("Failed to query workflow", zap.Error(err))
		//panic("Failed to query workflow.")
	}
	if err := resp.QueryResult.Get(&result); err != nil {
		common.LOCadenceHelper.Logger.Error("Failed to decode query result", zap.Error(err))
	}
	common.LOCadenceHelper.Logger.Info("Received consistent query result.", zap.Any("Result", result))
	return &result, err
}
*/
