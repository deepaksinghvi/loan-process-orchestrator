package loan_workflow

import (
	"context"
	"fmt"
	"github.com/deepaksinghvi/loan-process-orchestrator/common"
	"github.com/deepaksinghvi/loan-process-orchestrator/common/dto"
	"github.com/deepaksinghvi/loan-process-orchestrator/loan_workflow/loan_activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"strings"
	"time"
)

type LoanSignal struct {
	Message string // serializable
}

func LoanOriginationWorkflow(ctx workflow.Context, input dto.LoanApplicationInputStep) error {
	logger := workflow.GetLogger(ctx)
	currentState := "loan_app_submitted"
	err := workflow.SetQueryHandler(ctx, common.LoanWorkflowQueryType, func() (string, error) {
		return currentState, nil
	})
	if err != nil {
		currentState = "failed to register query handler"
		return err
	}
	loanApplicationOutput := dto.LoanApplicationOutputStep{}

	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        0, // unlimited retries
		NonRetryableErrorTypes: []string{"InvalidInputError", "ExternalIntegrationError"},
	}
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failed Activities by default.
		RetryPolicy: retrypolicy,
	}

	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)

	fmt.Printf("LoanOriginationWorkflow started")
	workflow.Sleep(ctx, time.Second*10) // workflow can be put on sleep
	/*
		Activities are invoked asynchronously through task lists.
		A task list is essentially a queue used to store an loan_activity task until it is picked up by an available worker.
	*/
	currentState = "loan_app_started"
	err = workflow.ExecuteActivity(ctx, loan_activity.LoanApplicationActivity, input).Get(ctx, &loanApplicationOutput)
	if err != nil {
		logger.Error("LoanApplicationActivity failed", err.Error())
		return err
	}

	creditDecisionInput := dto.CreditDecisionInputStep{
		ApplicationNo: loanApplicationOutput.ApplicationNo,
		AadhaarNumber: input.AadhaarNumber,
		PanNumber:     input.PanNumber,
		AccountNo:     loanApplicationOutput.AccountNo,
	}
	creditDecisionOutput := dto.CreditDecisionOutputStep{}
	currentState = "credit_check_initiated"
	err = workflow.ExecuteActivity(ctx, loan_activity.CreditDecisionInternalActivity, creditDecisionInput).Get(ctx, &creditDecisionOutput)
	if err != nil {
		logger.Error("CreditDecisionInternalActivity failed", err.Error())
		return err
	}
	var signalmessage string
	signalChan := workflow.GetSignalChannel(ctx, common.SignalName)
	signalChan.Receive(ctx, &signalmessage)
	loanFundingInput := dto.LoanFundingInputStep{
		ApplicationNo: creditDecisionOutput.ApplicationNo,
		AccountNo:     creditDecisionOutput.AccountNo,
	}
	loanFundingOutput := dto.LoanApplicationOutputStep{}
	currentState = "credit_check_completed_waiting_for_approval"
	if signalmessage == "APPROVED" {
		logger.Info("Signal Received as approved : " + signalmessage)
		currentState = "loan_approved"
		err = workflow.ExecuteActivity(ctx, loan_activity.LoanFundingActivity, loanFundingInput).Get(ctx, &loanFundingOutput)
		if err != nil {
			logger.Error("LoanFundingActivity failed", err.Error())
			return err
		}
	} else {
		logger.Info("Signal Received as rejected : " + signalmessage)
		currentState = "loan_rejected"
		err = workflow.ExecuteActivity(ctx, loan_activity.LoanRejectionActivity, loanFundingInput).Get(ctx, &loanFundingOutput)
		if err != nil {
			logger.Error("LoanFundingActivity failed", err.Error())
			return err
		}
	}

	logger.Info("Workflow result %v", loanFundingOutput)
	return nil
}

func checkConditionLoanApproved(ctx context.Context, signal string) (bool, error) {
	// some real logic happen here...
	return strings.Contains(signal, "APPROVED"), nil
}

func checkConditionLoanRejected(ctx context.Context, signal string) (bool, error) {
	// some real logic happen here...
	return strings.Contains(signal, "REJECTED"), nil
}

type conditionAndAction struct {
	// condition is a function pointer to a local loan_activity
	condition interface{}
	// action is a function pointer to a regulaloan_workflowr loan_activity
	action interface{}
}

var checks = []conditionAndAction{
	//{checkConditionLoanApproved, LoanFundingActivity},
	//{checkConditionLoanRejected, LoanRejectionActivity},
}
