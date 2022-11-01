package tg

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
	"go.opentelemetry.io/otel"
)

type Client struct {
	client *tgbotapi.BotAPI
	router IRouterTexte
}

type config interface {
	TelegramToken() string
}

type IRouterTexte interface {
	Execute(ctx context.Context, userID int64, textIn string, date time.Time) (string, error)
}

func New(cfg config, router IRouterTexte) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(cfg.TelegramToken())
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}

	return &Client{
		client: client,
		router: router,
	}, nil
}

func (c *Client) SendMessage(ctx context.Context, text string, userID int64) error {
	_, span := otel.Tracer("RatesUpdaterService").Start(ctx, "SendMessage")
	defer span.End()

	_, err := c.client.Send(tgbotapi.NewMessage(userID, text))
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}

	return nil
}

func (c *Client) Run(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if err := c.processing(ctx, update); err != nil {
				logger.Errorf("client processing error: %v", err)
			}
		case <-ctx.Done():
			c.client.StopReceivingUpdates()

			return
		}
	}
}

func (c *Client) processing(ctx context.Context, update tgbotapi.Update) error {
	ctx, span := otel.Tracer("tgClient").Start(ctx, "processing")
	defer span.End()

	if update.Message == nil {
		return nil
	}

	userID := update.Message.From.ID
	textIn := update.Message.Text
	date := time.Now()

	textOut, err := c.router.Execute(ctx, userID, textIn, date)
	if err != nil {
		return errors.Wrap(err, "client.processing")
	}

	logger.Infof("[%s][%d] %s -> %s", update.Message.From.UserName, userID, textIn, textOut)

	if len(textOut) == 0 {
		return nil
	}

	err = c.SendMessage(ctx, textOut, userID)
	if err != nil {
		return err
	}

	return nil
}
