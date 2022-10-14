package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/app"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/config"
)

func main() {
	configFile := config.ConfigFile

	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	cfg, err := config.New(configFile)
	if err != nil {
		log.Fatal("config init failed:", err)
	}

	app, err := app.New(cfg)
	if err != nil {
		log.Fatal("app init failed:", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	app.Run(ctx)

	log.Println("Exit")
}
