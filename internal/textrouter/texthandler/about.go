package texthandler

import (
	"context"
	"strings"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

type About struct{}

func NewAbout() *About {
	return &About{}
}

func (h *About) Name() string {
	return usecase.AboutCmdName
}

func (h *About) ConvertTextToCommand(ctx context.Context, text string, cmd *usecase.Command) bool {
	return strings.HasPrefix(text, "/about")
}

func (h *About) ConvertCommandToText(ctx context.Context, cmd *usecase.Command) (string, error) {
	return "Я бот для учета расходов. Автор @amyasnikov. OzonTech.", nil
}
