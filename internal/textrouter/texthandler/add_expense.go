package texthandler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

type IExpenseUsecaseAE interface {
	AddExpense(ctx context.Context, req usecase.AddExpenseReqDTO) error
}

type AddExpense struct {
	expenseUsecase IExpenseUsecaseAE
}

func NewAddExpense(expenseUsecase IExpenseUsecaseAE) *AddExpense {
	return &AddExpense{
		expenseUsecase: expenseUsecase,
	}
}

func (h *AddExpense) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *textrouter.Command) bool {
	categoryIndex := 1
	priceIndex := 2
	currencyIndex := 3
	argsCount := 4

	fields := strings.Fields(text)
	if len(fields) < 3 || len(fields) > argsCount || fields[0] != "расход" {
		return false
	}

	category := fields[categoryIndex]

	const bitSize = 64

	price, err := strconv.ParseFloat(fields[priceIndex], bitSize)
	if err != nil {
		return false
	}

	currency := ""
	if len(fields) > currencyIndex {
		currency = fields[currencyIndex]
	}

	cmd.AddExpenseReqDTO = &usecase.AddExpenseReqDTO{
		UserID:   userID,
		Category: category,
		Price:    price,
		Date:     date,
		Currency: currency,
	}

	return true
}

func (h *AddExpense) ExecuteCommand(ctx context.Context, cmd *textrouter.Command) error {
	if cmd.AddExpenseReqDTO == nil {
		return errors.Wrap(textrouter.ErrInvalidCommand, "AddExpense.ExecuteCommand")
	}

	// TODO добавить currency в ответ
	err := h.expenseUsecase.AddExpense(ctx, *cmd.AddExpenseReqDTO)
	if err != nil {
		return errors.Wrap(err, "AddExpense.ExecuteCommand")
	}

	return nil
}

func (h *AddExpense) ConvertCommandToText(cmd *textrouter.Command) (string, error) {
	if cmd.AddExpenseReqDTO == nil {
		return "", errors.Wrap(textrouter.ErrInvalidCommand, "AddExpense.ExecuteCommand")
	}

	textOut := fmt.Sprintf("Добавил %s - %0.2f %s %s", cmd.AddExpenseReqDTO.Category,
		cmd.AddExpenseReqDTO.Price, cmd.AddExpenseReqDTO.Currency,
		cmd.AddExpenseReqDTO.Date.Format(time.RFC1123))

	return textOut, nil
}
