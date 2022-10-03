package router

import (
	"time"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

type HandlerTextInterface interface {
	GetID() int
	ConvertTextToCommand(userID int64, text string, date time.Time, cmd *command)
	ExecuteCommand(cmd *command, productUsecase usecase.ProductUsecaseInterface) error
	ConvertCommandToText(cmd command) string
}
