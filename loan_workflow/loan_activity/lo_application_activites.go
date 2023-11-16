package loan_activity

import (
	"context"
	"fmt"
	"github.com/deepaksinghvi/loan-process-orchestrator/common/dto"
	"go.temporal.io/sdk/activity"
)

func LoanApplicationActivity(ctx context.Context, input dto.LoanApplicationInputStep) (dto.LoanApplicationOutputStep, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("LoanApplication started")
	output := dto.LoanApplicationOutputStep{
		ApplicationNo: activity.GetInfo(ctx).WorkflowExecution.RunID,
		AccountNo:     input.AccountNo,
	}
	logger.Info(fmt.Sprintf("LoanApplication for application no %s! for account no %s Completed.", output.ApplicationNo, output.AccountNo))
	return output, nil
}
