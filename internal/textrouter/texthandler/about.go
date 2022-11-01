package texthandler

import (
	"context"
	"strings"
	"time"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
)

type About struct{}

func NewAbout() *About {
	return &About{}
}

func (h *About) Name() string {
	return "about"
}

func (h *About) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *textrouter.Command) bool {
	return strings.HasPrefix(text, "/about")
}

func (h *About) ExecuteCommand(ctx context.Context, cmd *textrouter.Command) error {
	return nil
}

func (h *About) ConvertCommandToText(cmd *textrouter.Command) (string, error) {
	return "Я бот для учета расходов. Автор @amyasnikov. OzonTech.", nil
}
