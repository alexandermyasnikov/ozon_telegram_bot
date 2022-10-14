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
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/util"
)

type IExpenseUsecaseGR interface {
	GetReport(ctx context.Context, req usecase.GetReportReqDTO) (usecase.GetReportRespDTO, error)
}

type GetReport struct {
	expenseUsecase IExpenseUsecaseGR
}

func NewGetReport(expenseUsecase IExpenseUsecaseGR) *GetReport {
	return &GetReport{
		expenseUsecase: expenseUsecase,
	}
}

func (h *GetReport) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *textrouter.Command) bool {
	intervalIndex := 1
	currencyIndex := 2
	argsCount := 3

	fields := strings.Fields(text)
	if len(fields) == 0 || len(fields) > argsCount || fields[0] != "отчет" {
		return false
	}

	days := util.MonthToDays

	if len(fields) > intervalIndex {
		switch fields[intervalIndex] {
		case "день":
			days = util.DayToDays
		case "неделя":
			days = util.WeekToDays
		case "месяц":
			days = util.MonthToDays
		case "год":
			days = util.YearToDays
		default:
			return false
		}
	}

	currency := ""
	if len(fields) > currencyIndex {
		currency = fields[currencyIndex]
	}

	cmd.GetReportReqDTO = &usecase.GetReportReqDTO{
		UserID:   userID,
		Date:     date,
		Days:     days,
		Currency: currency,
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

	// TODO "дней" не склоняется
	textOut := fmt.Sprintf("Расходы по категориям за %d дней:\n", cmd.GetReportReqDTO.Days)

	lines := make([]string, 0, len(cmd.GetReportRespDTO.Categories))
	for k, v := range cmd.GetReportRespDTO.Categories {
		lines = append(lines, fmt.Sprintf("%s - %.2f", k, v.ToFloat()))
	}

	sort.Strings(lines)

	textOut += strings.Join(lines, "\n")

	return textOut, nil
}
