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

type ExpenseUsecaseAE interface {
	AddExpense(ctx context.Context, req usecase.AddExpenseReqDTO) (usecase.AddExpenseRespDTO, error)
}

type AddExpense struct {
	expenseUsecase ExpenseUsecaseAE
}

func NewAddExpense(expenseUsecase ExpenseUsecaseAE) *AddExpense {
	return &AddExpense{
		expenseUsecase: expenseUsecase,
	}
}

func (h *AddExpense) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *textrouter.Command) bool {
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
		UserID:   userID,
		Category: category,
		Price:    price,
		Date:     date,
	}

	return true
}

func (h *AddExpense) ExecuteCommand(ctx context.Context, cmd *textrouter.Command) error {
	if cmd.AddExpenseReqDTO == nil {
		return errors.Wrap(textrouter.ErrInvalidCommand, "AddExpense.ExecuteCommand")
	}

	resp, err := h.expenseUsecase.AddExpense(ctx, *cmd.AddExpenseReqDTO)
	if err != nil {
		return errors.Wrap(err, "AddExpense.ExecuteCommand")
	}

	cmd.AddExpenseRespDTO = &resp

	return nil
}

func (h *AddExpense) ConvertCommandToText(cmd *textrouter.Command) (string, error) {
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
