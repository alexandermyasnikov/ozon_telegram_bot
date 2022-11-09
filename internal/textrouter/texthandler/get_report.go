package texthandler

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/utils"
)

type GetReport struct{}

func NewGetReport() *GetReport {
	return &GetReport{}
}

func (h *GetReport) Name() string {
	return usecase.GetReportCmdName
}

func (h *GetReport) ConvertTextToCommand(ctx context.Context, text string, cmd *usecase.Command) bool {
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
		UserID:       cmd.UserID,
		Date:         cmd.Date,
		IntervalType: intervalType,
	}

	return true
}

func (h *GetReport) ConvertCommandToText(ctx context.Context, cmd *usecase.Command) (string, error) {
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
