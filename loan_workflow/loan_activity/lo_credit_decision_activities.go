package loan_activity

import (
	"context"
	"fmt"
	"github.com/deepaksinghvi/loan-process-orchestrator/common/dto"
	"github.com/pborman/uuid"
	"go.temporal.io/sdk/activity"
	"time"
)

func CreditDecisionInternalActivity(ctx context.Context, input dto.CreditDecisionInputStep) (dto.CreditDecisionOutputStep, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("CreditDecisionInternalActivity started")
	output := dto.CreditDecisionOutputStep{
		ApplicationNo: input.ApplicationNo,
		AccountNo:     uuid.New(),
		CreditScore:   800.0,
	}
	logger.Info(fmt.Sprintf("CreditDecisionInternalActivity for application no %s! with credit score %f Completed.", input.ApplicationNo, output.CreditScore))
	return output, nil
}

func CreditDecisionExternalActivity(ctx context.Context, input dto.CreditDecisionInputStep) (dto.CreditDecisionOutputStep, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("CreditDecisionInternalActivity started")
	output := dto.CreditDecisionOutputStep{
		ApplicationNo: input.ApplicationNo,
		AccountNo:     uuid.New(),
		CreditScore:   800.0,
	}
	for i := 0; i < 5; i++ {
		time.Sleep(60 * time.Second)
		activity.RecordHeartbeat(ctx, i+1)
	}
	logger.Info(fmt.Sprintf("CreditDecisionInternalActivity for application no %s! with credit score %f Completed.", input.ApplicationNo, output.CreditScore))
	return output, nil
}
