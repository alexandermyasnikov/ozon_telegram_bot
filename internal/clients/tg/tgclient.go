package tg

import (
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/router"
)

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

func (c *Client) SendMessage(text string, userID int64) error {
	_, err := c.client.Send(tgbotapi.NewMessage(userID, text))
	if err != nil {
		return errors.Wrap(err, "client.Send")
	}

	return nil
}

func (c *Client) ListenUpdates(router router.RouterTextInterface) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.client.GetUpdatesChan(u)

	log.Println("listening for messages")

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		textIn := update.Message.Text
		date := time.Now()

		textOut, err := router.Execute(userID, textIn, date)
		if err != nil {
			log.Println("error processing message:", err)

			continue
		}

		log.Printf("[%s][%d] %s -> %s", update.Message.From.UserName, userID, textIn, textOut)

		if len(textOut) == 0 {
			continue
		}

		err = c.SendMessage(textOut, userID)
		if err != nil {
			log.Println("error processing message:", err)
		}
	}
}
