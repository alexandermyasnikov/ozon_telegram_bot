package tg

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
	"go.opentelemetry.io/otel"
)

type MsgCallback = func(ctx context.Context, userID int64, date time.Time, text string)

type Client struct {
	client *tgbotapi.BotAPI
}

type config interface {
	TelegramToken() string
}

func New(cfg config) (*Client, error) {
	client, err := tgbotapi.NewBotAPI(cfg.TelegramToken())
	if err != nil {
		return nil, errors.Wrap(err, "NewBotAPI")
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Write(ctx context.Context, text string, userID int64) error {
	_, span := otel.Tracer("tgClient").Start(ctx, "WriteMessage")
	defer span.End()

	logger.Infof("client.Write [%d][%s]", userID, text)

	_, err := c.client.Send(tgbotapi.NewMessage(userID, text))
	if err != nil {
		return errors.Wrap(err, "client.Write")
	}

	return nil
}

func (c *Client) Read(ctx context.Context, callback MsgCallback) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			c.processing(ctx, update, callback)
		case <-ctx.Done():
			c.client.StopReceivingUpdates()

			return
		}
	}
}

func (c *Client) processing(ctx context.Context, update tgbotapi.Update, callback MsgCallback) {
	ctx, span := otel.Tracer("tgClient").Start(ctx, "processing")
	defer span.End()

	if update.Message == nil {
		return
	}

	logger.Infof("client.read: [%s][%d][%s]", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)

	callback(ctx, update.Message.From.ID, update.Message.Time(), update.Message.Text)
}
