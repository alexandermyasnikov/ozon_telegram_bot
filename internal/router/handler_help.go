package router

import (
	"strings"
	"time"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

var _ HandlerTextInterface = (*HandlerTextHelp)(nil)

type HandlerTextHelp struct{}

func (h *HandlerTextHelp) GetID() int {
	return cmdHelp
}

func (h *HandlerTextHelp) ConvertTextToCommand(userID int64, text string, date time.Time, cmd *command) {
	if strings.HasPrefix(text, "/help") {
		cmd.id = h.GetID()
	}
}

func (h *HandlerTextHelp) ExecuteCommand(cmd *command, productUsecase usecase.ProductUsecaseInterface) error {
	return nil
}

func (h *HandlerTextHelp) ConvertCommandToText(cmd command) string {
	return cmdHelpText
}
