package app

import (
	"log"

	productStorage "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/router"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
)

type client interface {
	ListenUpdates(router router.RouterTextInterface)
	SendMessage(text string, userID int64) error
}

type config interface {
	TelegramToken() string
}

var _ client = &tg.Client{}

type App struct {
	client client
	router router.RouterTextInterface
}

func New(cfg config) (App, error) {
	tgClient, err := tg.New(cfg)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	productStorage := productStorage.NewProductStorage()
	productUsecase := usecase.NewProductUsecase(productStorage)

	routerText := router.NewRouterText(productUsecase)

	routerText.Register(&router.HandlerTextStart{})
	routerText.Register(&router.HandlerTextUnknown{})
	routerText.Register(&router.HandlerTextHelp{})
	routerText.Register(&router.HandlerTextAbout{})
	routerText.Register(&router.HandlerTextAddProduct{})
	routerText.Register(&router.HandlerTextGetStats{})

	return App{
		client: tgClient,
		router: routerText,
	}, nil
}

func (a *App) Run() {
	a.client.ListenUpdates(a.router)
}
