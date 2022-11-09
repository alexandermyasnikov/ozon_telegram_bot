package texthandler

import (
	"context"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

type Unknown struct{}

func NewUnknown() *Unknown {
	return &Unknown{}
}

func (h *Unknown) Name() string {
	return usecase.UnknownCmdName
}

func (h *Unknown) ConvertTextToCommand(ctx context.Context, text string, cmd *usecase.Command) bool {
	return true
}

func (h *Unknown) ConvertCommandToText(ctx context.Context, cmd *usecase.Command) (string, error) {
	return "Не могу понять. Введи /help для более подробной информации", nil
}
