package apptgclientwriter

import (
	"context"
	"encoding/json"

	kafkareader "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/kafka/kafka_reader"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/config"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter/texthandler"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
)

type AppTgClientWriter struct {
	reader   *kafkareader.KafkaReader
	callback kafkareader.MsgCallback
}

type Client interface {
	Write(context.Context, string, int64) error
}

func New(ctx context.Context, cfg *config.Config) (AppTgClientWriter, error) {
	tgClient, err := tg.New(cfg)
	if err != nil {
		logger.Fatalf("tg client init failed: %v", err)
	}

	return NewWithCustomClient(ctx, cfg, tgClient)
}

func NewWithCustomClient(ctx context.Context, cfg *config.Config, client Client) (AppTgClientWriter, error) {
	reader := kafkareader.New(cfg.GetKafkaAddr(), usecase.ProcessCmdState, "tgClientReader")

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

	callback := func(ctx context.Context, key, value []byte) {
		var cmd usecase.Command

		err := json.Unmarshal(value, &cmd)
		if err != nil {
			logger.Errorf("can not unmarshal command: %v", err)
		}

		text := routerText.ConvertCommandToText(ctx, &cmd)

		err = client.Write(ctx, text, cmd.UserID)
		if err != nil {
			logger.Errorf("can not write message: %v", err)
		}
	}

	return AppTgClientWriter{
		reader:   reader,
		callback: callback,
	}, nil
}

func (a *AppTgClientWriter) Run(ctx context.Context) {
	a.reader.Read(ctx, a.callback)
}
