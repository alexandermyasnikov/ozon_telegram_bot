package router

import (
	"strings"
	"time"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

var _ HandlerTextInterface = (*HandlerTextStart)(nil)

type HandlerTextStart struct{}

func (h *HandlerTextStart) GetID() int {
	return cmdStart
}

func (h *HandlerTextStart) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *command) {
	if strings.HasPrefix(text, "/start") {
		cmd.id = h.GetID()
	}
}

func (h *HandlerTextStart) ExecuteCommand(cmd *command, productUsecase usecase.ProductUsecaseInterface) error {
	return nil
}

func (h *HandlerTextStart) ConvertCommandToText(cmd command) string {
	return cmdStartText
}
