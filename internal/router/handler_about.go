package router

import (
	"strings"
	"time"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

var _ HandlerTextInterface = (*HandlerTextAbout)(nil)

type HandlerTextAbout struct{}

func (h *HandlerTextAbout) GetID() int {
	return cmdAbout
}

func (h *HandlerTextAbout) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *command) {
	if strings.HasPrefix(text, "/about") {
		cmd.id = h.GetID()
	}
}

func (h *HandlerTextAbout) ExecuteCommand(cmd *command, productUsecase usecase.ProductUsecaseInterface) error {
	return nil
}

func (h *HandlerTextAbout) ConvertCommandToText(cmd command) string {
	return cmdAboutText
}
