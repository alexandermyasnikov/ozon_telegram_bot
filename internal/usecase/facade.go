package usecase

import (
	"context"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
	"go.opentelemetry.io/otel"
)

type FacadeUsecase struct {
	expenseUsecase ExpenseUsecase
}

func New(expenseUsecase *ExpenseUsecase) *FacadeUsecase {
	return &FacadeUsecase{
		expenseUsecase: *expenseUsecase,
	}
}

func forward[Req, Resp any](ctx context.Context, usecase func(context.Context, Req) (Resp, error),
	req *Req, resp **Resp,
) error {
	if req == nil {
		return errors.New("internal error")
	}

	r, err := usecase(ctx, *req)
	*resp = &r

	return err
}

func (f *FacadeUsecase) ExecuteCommand(ctx context.Context, cmd *Command) error {
	logger.Infof("FacadeUsecase.ExecuteCommand: %v", cmd.Name)

	ctx, span := otel.Tracer("FacadeUsecase").Start(ctx, "ExecuteCommand")
	defer span.End()

	switch cmd.Name {
	case SetCurrencyCmdName:
		return forward(ctx, f.expenseUsecase.SetDefaultCurrency, cmd.SetDefaultCurrencyReqDTO, &cmd.SetDefaultCurrencyRespDTO)
	case AddExpenseCmdName:
		return forward(ctx, f.expenseUsecase.AddExpense, cmd.AddExpenseReqDTO, &cmd.AddExpenseRespDTO)
	case GetReportCmdName:
		return forward(ctx, f.expenseUsecase.GetReport, cmd.GetReportReqDTO, &cmd.GetReportRespDTO)
	case SetLimitCmdName:
		return forward(ctx, f.expenseUsecase.SetLimit, cmd.SetLimitReqDTO, &cmd.SetLimitRespDTO)
	case GetLimitsCmdName:
		return forward(ctx, f.expenseUsecase.GetLimits, cmd.GetLimitsReqDTO, &cmd.GetLimitsRespDTO)
	case StartCmdName:
	case HelpCmdName:
	case AboutCmdName:
	case UnknownCmdName:
	default:
		logger.Errorf("unknown command: %v", cmd.Name)
	}

	return nil
}
