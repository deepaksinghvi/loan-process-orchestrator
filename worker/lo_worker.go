package worker

import (
	"context"
	"github.com/deepaksinghvi/loan-process-orchestrator/common"
	"github.com/deepaksinghvi/loan-process-orchestrator/common/dto"
	"github.com/deepaksinghvi/loan-process-orchestrator/loan_workflow"
	"github.com/deepaksinghvi/loan-process-orchestrator/loan_workflow/loan_activity"
	"github.com/google/uuid"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"log"
)

func StartWorrker() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, common.LOTaskQueueName, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterWorkflowWithOptions(loan_workflow.LoanOriginationWorkflow, workflow.RegisterOptions{Name: "loan_workflow.LoanOriginationWorkflow"})
	w.RegisterActivityWithOptions(loan_activity.LoanApplicationActivity, activity.RegisterOptions{Name: "loan_activity.LoanApplicationActivity"})
	w.RegisterActivityWithOptions(loan_activity.CreditDecisionInternalActivity, activity.RegisterOptions{Name: "loan_activity.CreditDecisionInternalActivity"})
	w.RegisterActivityWithOptions(loan_activity.LoanFundingActivity, activity.RegisterOptions{Name: "loan_activity.LoanFundingActivity"})
	w.RegisterActivityWithOptions(loan_activity.LoanRejectionActivity, activity.RegisterOptions{Name: "loan_activity.LoanRejectionActivity"})

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func StartLoanWorkflowExecution(input dto.LoanApplicationInputStep) (string, string, error) {

	// Create the client object just once per process
	c, err := client.Dial(client.Options{})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
		return "", "", err
	}

	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: common.LOTaskQueueName,
	}

	log.Printf("Starting loan application for Applicant: %s with  AccountNo: %s", input.ApplicantName, input.AccountNo)

	we, err := c.ExecuteWorkflow(context.Background(), options, loan_workflow.LoanOriginationWorkflow, input)
	if err != nil {
		log.Fatalln("Unable to start the Workflow:", err)
		return "", "", err
	}

	log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())
	return we.GetID(), we.GetRunID(), nil
}

func SendSignalToWorkflow(workflowID, runID, signalName, signalValue string) error {

	// Create the client object just once per process
	c, err := client.Dial(client.Options{})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
		return err
	}

	defer c.Close()
	err = c.SignalWorkflow(context.Background(), workflowID, runID, signalName, signalValue)
	if err != nil {
		log.Printf("Unable to send singla to Workflow: %s with error: %s", workflowID, err.Error())
		return err
	}
	return nil
}

func QueryWorkflow(workflowID, runID, queryType string) (string, error) {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
		return "", err
	}

	defer c.Close()
	response, err := c.QueryWorkflow(context.Background(), workflowID, runID, queryType)
	if err != nil {
		log.Printf("Unable to send singla to Workflow: %s with error: %s", workflowID, err.Error())
		return "", err
	}
	var result string
	err = response.Get(&result)
	if err != nil {
		log.Printf("Unable to send singla to Workflow: %s with error: %s", workflowID, err.Error())
		return "", err
	}
	log.Println("Received Query result. Result: " + result)
	return result, nil
}
