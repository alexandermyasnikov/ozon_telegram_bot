package router

import (
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

type RouterTextInterface interface {
	Execute(userID int64, textIn string, date time.Time) (string, error)
	Register(handler HandlerTextInterface)
}

type RouterText struct {
	productUsecase usecase.ProductUsecaseInterface
	handlers       map[int]HandlerTextInterface
}

var _ RouterTextInterface = (*RouterText)(nil)

func NewRouterText(productUsecase usecase.ProductUsecaseInterface) *RouterText {
	return &RouterText{
		productUsecase: productUsecase,
		handlers:       make(map[int]HandlerTextInterface),
	}
}

func (r *RouterText) Register(handler HandlerTextInterface) {
	r.handlers[handler.GetID()] = handler
}

func (r *RouterText) Execute(userID int64, textIn string, date time.Time) (string, error) {
	if len(textIn) == 0 {
		return "", nil
	}

	command := r.convertTextToCommand(userID, textIn, date)

	err := r.executeCommand(&command)
	if err != nil {
		return "", err
	}

	textOut := r.convertCommandToText(command)

	return textOut, nil
}

func (r *RouterText) convertTextToCommand(userID int64, text string, date time.Time) command {
	cmd := command{
		id:                   cmdUnknown,
		addProductReqDTO:     nil,
		getStatisticsReqDTO:  nil,
		getStatisticsRespDTO: nil,
	}

	for _, handler := range r.handlers {
		handler.ConvertTextToCommand(userID, text, date, &cmd)

		if cmd.id != cmdUnknown {
			break
		}
	}

	return cmd
}

func (r *RouterText) executeCommand(cmd *command) error {
	handler, ok := r.handlers[cmd.id]
	if !ok {
		return nil
	}

	err := handler.ExecuteCommand(cmd, r.productUsecase)

	return errors.Wrap(err, "handler.ExecuteCommand")
}

func (r *RouterText) convertCommandToText(cmd command) string {
	handler, ok := r.handlers[cmd.id]
	if !ok {
		return ""
	}

	return handler.ConvertCommandToText(cmd)
}
