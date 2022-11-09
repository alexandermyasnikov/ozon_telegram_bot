package texthandler

import (
	"context"
	"strings"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

type Start struct{}

func NewStart() *Start {
	return &Start{}
}

func (h *Start) Name() string {
	return usecase.StartCmdName
}

func (h *Start) ConvertTextToCommand(ctx context.Context, text string, cmd *usecase.Command) bool {
	return strings.HasPrefix(text, "/start")
}

func (h *Start) ConvertCommandToText(ctx context.Context, cmd *usecase.Command) (string, error) {
	return `Привет. Напиши свои расходы и я запомпю их.
Введи /help для более подробной информации`, nil
}
