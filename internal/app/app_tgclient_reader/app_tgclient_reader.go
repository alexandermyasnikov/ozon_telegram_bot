package apptgclientreader

import (
	"context"
	"encoding/json"
	"time"

	kafkawriter "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/kafka/kafka_writer"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/config"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter/texthandler"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
)

type AppTgClientReader struct {
	client   Client
	callback tg.MsgCallback
}

type Client interface {
	Read(context.Context, func(context.Context, int64, time.Time, string))
}

func New(ctx context.Context, cfg *config.Config) (AppTgClientReader, error) {
	tgClient, err := tg.New(cfg)
	if err != nil {
		logger.Fatalf("tg client init failed: %v", err)
	}

	return NewWithCustomClient(ctx, cfg, tgClient)
}

func NewWithCustomClient(ctx context.Context, cfg *config.Config, client Client) (AppTgClientReader, error) {
	writer := kafkawriter.New(cfg.GetKafkaAddr(), usecase.ReadCmdState)

	routerText := textrouter.New()

	routerText.Register(texthandler.NewStart())
	routerText.Register(texthandler.NewHelp())
	routerText.Register(texthandler.NewAbout())
	routerText.Register(texthandler.NewSetDefaultCurrency())
	routerText.Register(texthandler.NewAddExpense())
	routerText.Register(texthandler.NewGetReport())
	routerText.Register(texthandler.NewSetLimit())
	routerText.Register(texthandler.NewGetLimits())
	routerText.Register(texthandler.NewUnknown())

	callback := func(ctx context.Context, userID int64, date time.Time, text string) {
		cmd := routerText.ConvertTextToCommand(ctx, userID, date, text)

		buf, err := json.Marshal(cmd)
		if err != nil {
			logger.Errorf("can not marshal command: %v", err)
		}

		err = writer.Write(ctx, []byte(cmd.Name), buf)
		if err != nil {
			logger.Errorf("can not write message: %v", err)
		}
	}

	return AppTgClientReader{
		client:   client,
		callback: callback,
	}, nil
}

func (a *AppTgClientReader) Run(ctx context.Context) {
	a.client.Read(ctx, a.callback)
}
