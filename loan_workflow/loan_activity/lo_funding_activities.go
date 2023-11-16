package loan_activity

import (
	"context"
	"fmt"
	"github.com/deepaksinghvi/loan-process-orchestrator/common/dto"
	"go.temporal.io/sdk/activity"
)

func LoanFundingActivity(ctx context.Context, input dto.LoanFundingInputStep) (dto.LoanFundingOutputStep, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("LoanFundingActivity started")
	output := dto.LoanFundingOutputStep{
		ApplicationNo:      input.ApplicationNo,
		AccountNo:          input.AccountNo,
		DisbursementAmount: 1000.0,
	}
	logger.Info(fmt.Sprintf("LoanFundingActivity for application no %s! with disbursement amount %f Completed.", input.ApplicationNo, output.DisbursementAmount))
	return output, nil
}

func LoanRejectionActivity(ctx context.Context, input dto.LoanFundingInputStep) (dto.LoanFundingOutputStep, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("LoanRejectionActivity started")
	output := dto.LoanFundingOutputStep{
		ApplicationNo:      input.ApplicationNo,
		AccountNo:          input.AccountNo,
		DisbursementAmount: 0.0,
	}
	logger.Info(fmt.Sprintf("LoanRejectionActivity for application no %s! with Rejection of application Completed.", input.ApplicationNo))
	return output, nil
}
