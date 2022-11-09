package texthandler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/utils"
)

type AddExpense struct{}

func NewAddExpense() *AddExpense {
	return &AddExpense{}
}

func (h *AddExpense) Name() string {
	return usecase.AddExpenseCmdName
}

func (h *AddExpense) ConvertTextToCommand(ctx context.Context, text string, cmd *usecase.Command) bool {
	categoryIndex := 1
	priceIndex := 2
	argsCountMin := 3
	argsCountMax := 3

	fields := strings.Fields(text)
	if len(fields) < argsCountMin || len(fields) > argsCountMax || fields[0] != "расход" {
		return false
	}

	category := fields[categoryIndex]

	price, err := decimal.NewFromString(fields[priceIndex])
	if err != nil {
		return false
	}

	cmd.AddExpenseReqDTO = &usecase.AddExpenseReqDTO{
		UserID:   cmd.UserID,
		Category: category,
		Price:    price,
		Date:     cmd.Date,
	}

	return true
}

func (h *AddExpense) ConvertCommandToText(ctx context.Context, cmd *usecase.Command) (string, error) {
	if cmd.AddExpenseReqDTO == nil || cmd.AddExpenseRespDTO == nil {
		return "", errors.Wrap(textrouter.ErrInvalidCommand, "AddExpense.ExecuteCommand")
	}

	precision := 2

	textOut := fmt.Sprintf("Добавил %s - %s %s %s", cmd.AddExpenseReqDTO.Category,
		cmd.AddExpenseReqDTO.Price.StringFixed(int32(precision)), cmd.AddExpenseRespDTO.Currency,
		cmd.AddExpenseReqDTO.Date.Format(time.RFC1123))

	for _, interval := range []int{utils.DayInterval, utils.WeekInterval, utils.MonthInterval} {
		limit, ok := cmd.AddExpenseRespDTO.Limits[interval]
		if !ok {
			continue
		}

		if limit.GreaterThanOrEqual(decimal.Zero) {
			continue
		}

		intervalStr, _ := utils.IntervalToStr(interval)

		textOut += fmt.Sprintf("\nВнимание! Превышен лимит: %s - %s",
			intervalStr, limit.Neg().StringFixed(int32(precision)))
	}

	return textOut, nil
}
