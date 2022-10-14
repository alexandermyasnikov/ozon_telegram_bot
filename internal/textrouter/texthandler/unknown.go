package texthandler

import (
	"context"
	"time"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
)

type Unknown struct{}

func NewUnknown() *Unknown {
	return &Unknown{}
}

func (h *Unknown) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *textrouter.Command) bool {
	return true
}

func (h *Unknown) ExecuteCommand(ctx context.Context, cmd *textrouter.Command) error {
	return nil
}

func (h *Unknown) ConvertCommandToText(cmd *textrouter.Command) (string, error) {
	return "Не могу понять. Введи /help для более подробной информации", nil
}
