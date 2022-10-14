package texthandler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

type IExpenseUsecaseSDC interface {
	SetDefaultCurrency(context.Context, usecase.SetDefaultCurrencyReqDTO) error
}

type SetDefaultCurrency struct {
	expenseUsecase IExpenseUsecaseSDC
}

func NewSetDefaultCurrency(expenseUsecase IExpenseUsecaseSDC) *SetDefaultCurrency {
	return &SetDefaultCurrency{
		expenseUsecase: expenseUsecase,
	}
}

func (h *SetDefaultCurrency) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *textrouter.Command,
) bool {
	currencyIndex := 1
	argsCount := 2

	fields := strings.Fields(text)
	if len(fields) != argsCount || fields[0] != "валюта" {
		return false
	}

	currency := fields[currencyIndex]

	cmd.SetDefaultCurrencyReqDTO = &usecase.SetDefaultCurrencyReqDTO{
		UserID:   userID,
		Currency: currency,
	}

	return true
}

func (h *SetDefaultCurrency) ExecuteCommand(ctx context.Context, cmd *textrouter.Command) error {
	if cmd.SetDefaultCurrencyReqDTO == nil {
		return errors.Wrap(textrouter.ErrInvalidCommand, "SetDefaultCurrency.ExecuteCommand")
	}

	err := h.expenseUsecase.SetDefaultCurrency(ctx, *cmd.SetDefaultCurrencyReqDTO)

	return errors.Wrap(err, "SetDefaultCurrency.ExecuteCommand")
}

func (h *SetDefaultCurrency) ConvertCommandToText(cmd *textrouter.Command) (string, error) {
	if cmd.SetDefaultCurrencyReqDTO == nil {
		return "", errors.Wrap(textrouter.ErrInvalidCommand, "SetDefaultCurrency.ExecuteCommand")
	}

	textOut := fmt.Sprintf("Задана валюта по умолчанию %s", cmd.SetDefaultCurrencyReqDTO.Currency)

	return textOut, nil
}
