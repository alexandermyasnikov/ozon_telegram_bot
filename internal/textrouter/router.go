package textrouter

import (
	"context"
	"time"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
)

type Handler interface {
	Name() string
	ConvertTextToCommand(ctx context.Context, text string, cmd *usecase.Command) bool
	ConvertCommandToText(ctx context.Context, cmd *usecase.Command) (string, error)
}

type RouterText struct {
	handlers []Handler
}

func New() *RouterText {
	return &RouterText{
		handlers: make([]Handler, 0),
	}
}

func (r *RouterText) Register(handler Handler) {
	r.handlers = append(r.handlers, handler)
}

func (r *RouterText) ConvertTextToCommand(ctx context.Context, userID int64, date time.Time, text string,
) usecase.Command {
	cmd := usecase.Command{
		MessageInfo: usecase.MessageInfo{
			UserID: userID,
			Date:   date,
		},
	}

	for _, handler := range r.handlers {
		if handler.ConvertTextToCommand(ctx, text, &cmd) {
			cmd.Name = handler.Name()

			break
		}
	}

	return cmd
}

func (r *RouterText) ConvertCommandToText(ctx context.Context, cmd *usecase.Command) string {
	for _, handler := range r.handlers {
		if handler.Name() == cmd.Name {
			text, err := handler.ConvertCommandToText(ctx, cmd)
			if err != nil {
				logger.Errorf("can not convert command to text: %v", err)

				return ErrInvalidCommand.Error()
			}

			return text
		}
	}

	logger.Errorf("unknown handler: %v", cmd.Name)

	return ErrInvalidCommand.Error()
}
