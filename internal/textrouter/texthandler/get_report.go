package texthandler

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/utils"
)

type ExpenseUsecaseGR interface {
	GetReport(ctx context.Context, req usecase.GetReportReqDTO) (usecase.GetReportRespDTO, error)
}

type GetReport struct {
	expenseUsecase ExpenseUsecaseGR
}

func NewGetReport(expenseUsecase ExpenseUsecaseGR) *GetReport {
	return &GetReport{
		expenseUsecase: expenseUsecase,
	}
}

func (h *GetReport) Name() string {
	return "getReport"
}

func (h *GetReport) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *textrouter.Command) bool {
	intervalIndex := 1
	argsCountMin := 2
	argsCountMax := 2

	fields := strings.Fields(text)
	if len(fields) < argsCountMin || len(fields) > argsCountMax || fields[0] != "отчет" {
		return false
	}

	var intervalType int

	if len(fields) > intervalIndex {
		interval, ok := utils.IntervalFromStr(fields[intervalIndex])

		if !ok {
			return false
		}

		intervalType = interval
	}

	cmd.GetReportReqDTO = &usecase.GetReportReqDTO{
		UserID:       userID,
		Date:         date,
		IntervalType: intervalType,
	}

	return true
}

func (h *GetReport) ExecuteCommand(ctx context.Context, cmd *textrouter.Command) error {
	if cmd.GetReportReqDTO == nil {
		return errors.Wrap(textrouter.ErrInvalidCommand, "GetReport.ExecuteCommand")
	}

	resp, err := h.expenseUsecase.GetReport(ctx, *cmd.GetReportReqDTO)
	if err != nil {
		return errors.Wrap(err, "GetReport.ExecuteCommand")
	}

	cmd.GetReportRespDTO = &resp

	return nil
}

func (h *GetReport) ConvertCommandToText(cmd *textrouter.Command) (string, error) {
	if cmd.GetReportReqDTO == nil || cmd.GetReportRespDTO == nil {
		return "", errors.Wrap(textrouter.ErrInvalidCommand, "GetReport.ExecuteCommand")
	}

	precision := 2

	intervalType, _ := utils.IntervalToStr(cmd.GetReportReqDTO.IntervalType)

	textOut := fmt.Sprintf("Расходы по категориям за %s:\n", intervalType)

	lines := make([]string, 0, len(cmd.GetReportRespDTO.Expenses))
	for _, expense := range cmd.GetReportRespDTO.Expenses {
		lines = append(lines, fmt.Sprintf("%s - %s", expense.Category, expense.Sum.StringFixed(int32(precision))))
	}

	sort.Strings(lines)

	textOut += strings.Join(lines, "\n")

	return textOut, nil
}
