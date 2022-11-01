package textrouter

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/metrics"
	"go.opentelemetry.io/otel"
)

type Command struct {
	SetDefaultCurrencyReqDTO *usecase.SetDefaultCurrencyReqDTO
	AddExpenseReqDTO         *usecase.AddExpenseReqDTO
	AddExpenseRespDTO        *usecase.AddExpenseRespDTO
	GetReportReqDTO          *usecase.GetReportReqDTO
	GetReportRespDTO         *usecase.GetReportRespDTO
	SetLimitReqDTO           *usecase.SetLimitReqDTO
	SetLimitRespDTO          *usecase.SetLimitRespDTO
	GetLimitsReqDTO          *usecase.GetLimitsReqDTO
	GetLimitsRespDTO         *usecase.GetLimitsRespDTO
}

type Handler interface {
	Name() string
	ConvertTextToCommand(userID int64, text string, date time.Time, cmd *Command) bool
	ExecuteCommand(context.Context, *Command) error
	ConvertCommandToText(cmd *Command) (string, error)
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

func (r *RouterText) Execute(ctx context.Context, userID int64, textIn string, date time.Time) (string, error) {
	ctx, span := otel.Tracer("RouterText").Start(ctx, "Execute")
	defer span.End()

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

	err := r.executeCommand(ctx, handler, &cmd)
	if err != nil {
		return "", errors.Wrap(err, "handler.ExecuteCommand")
	}

	textOut, err := handler.ConvertCommandToText(&cmd)

	return textOut, errors.Wrap(err, "RouterText.Executer")
}

func (r *RouterText) executeCommand(ctx context.Context, handler Handler, cmd *Command) error {
	ctx, span := otel.Tracer("RouterText").Start(ctx, "ExecuteCommand")
	defer span.End()

	startTime := time.Now()
	err := handler.ExecuteCommand(ctx, cmd)
	duration := time.Since(startTime)

	metrics.SummaryExecuteTimeObserve(handler.Name(), duration.Seconds())
	metrics.CounterMsgInc(handler.Name())

	return errors.Wrap(err, "RouterText.executeCommand")
}
