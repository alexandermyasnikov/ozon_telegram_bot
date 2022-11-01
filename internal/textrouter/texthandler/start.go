package texthandler

import (
	"context"
	"strings"
	"time"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
)

type Start struct{}

func NewStart() *Start {
	return &Start{}
}

func (h *Start) Name() string {
	return "start"
}

func (h *Start) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *textrouter.Command) bool {
	return strings.HasPrefix(text, "/start")
}

func (h *Start) ExecuteCommand(ctx context.Context, cmd *textrouter.Command) error {
	return nil
}

func (h *Start) ConvertCommandToText(cmd *textrouter.Command) (string, error) {
	return `Привет. Напиши свои расходы и я запомпю их.
Введи /help для более подробной информации`, nil
}
