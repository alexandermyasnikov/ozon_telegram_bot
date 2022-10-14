package textrouter

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

type Command struct {
	// TODO нужно обобщить, в тестах громоздко все поля заполнять
	SetDefaultCurrencyReqDTO *usecase.SetDefaultCurrencyReqDTO
	AddExpenseReqDTO         *usecase.AddExpenseReqDTO
	GetReportReqDTO          *usecase.GetReportReqDTO
	GetReportRespDTO         *usecase.GetReportRespDTO
}

type IHandler interface {
	ConvertTextToCommand(userID int64, text string, date time.Time, cmd *Command) bool
	ExecuteCommand(context.Context, *Command) error
	ConvertCommandToText(cmd *Command) (string, error)
}

type RouterText struct {
	handlers []IHandler
}

func New() *RouterText {
	return &RouterText{
		handlers: make([]IHandler, 0),
	}
}

func (r *RouterText) Register(handler IHandler) {
	r.handlers = append(r.handlers, handler)
}

func (r *RouterText) Execute(ctx context.Context, userID int64, textIn string, date time.Time) (string, error) {
	if len(textIn) == 0 {
		return "", nil
	}

	var cmd Command

	handlerID := -1

	for id, handler := range r.handlers {
		if handler.ConvertTextToCommand(userID, textIn, date, &cmd) {
			handlerID = id

			break
		}
	}

	if handlerID == -1 {
		return "", nil
	}

	handler := r.handlers[handlerID]

	err := handler.ExecuteCommand(ctx, &cmd)
	if err != nil {
		return "", errors.Wrap(err, "handler.ExecuteCommand")
	}

	textOut, err := handler.ConvertCommandToText(&cmd)

	return textOut, errors.Wrap(err, "RouterText.Executer")
}
