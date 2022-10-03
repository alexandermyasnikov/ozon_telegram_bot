package router

import (
	"time"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

var _ HandlerTextInterface = (*HandlerTextUnknown)(nil)

type HandlerTextUnknown struct{}

func (h *HandlerTextUnknown) GetID() int {
	return cmdUnknown
}

func (h *HandlerTextUnknown) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *command) {
}

func (h *HandlerTextUnknown) ExecuteCommand(cmd *command, productUsecase usecase.ProductUsecaseInterface) error {
	return nil
}

func (h *HandlerTextUnknown) ConvertCommandToText(cmd command) string {
	return cmdUnknownText
}
