package texthandler

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

type SetDefaultCurrency struct{}

func NewSetDefaultCurrency() *SetDefaultCurrency {
	return &SetDefaultCurrency{}
}

func (h *SetDefaultCurrency) Name() string {
	return usecase.SetCurrencyCmdName
}

func (h *SetDefaultCurrency) ConvertTextToCommand(ctx context.Context, text string, cmd *usecase.Command) bool {
	currencyIndex := 1
	argsCount := 2

	fields := strings.Fields(text)
	if len(fields) != argsCount || fields[0] != "валюта" {
		return false
	}

	currency := fields[currencyIndex]

	cmd.SetDefaultCurrencyReqDTO = &usecase.SetDefaultCurrencyReqDTO{
		UserID:   cmd.UserID,
		Currency: currency,
	}

	return true
}

func (h *SetDefaultCurrency) ConvertCommandToText(ctx context.Context, cmd *usecase.Command) (string, error) {
	if cmd.SetDefaultCurrencyReqDTO == nil {
		return "", errors.Wrap(textrouter.ErrInvalidCommand, "SetDefaultCurrency.ExecuteCommand")
	}

	textOut := fmt.Sprintf("Задана валюта по умолчанию %s", cmd.SetDefaultCurrencyReqDTO.Currency)

	return textOut, nil
}
