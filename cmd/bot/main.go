package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/app"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/config"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
)

func main() {
	var configPath = flag.String("config", config.ConfigFile, "path to config file")

	flag.Parse()

	cfg, err := config.New(*configPath)
	if err != nil {
		logger.Fatalf("config init failed: %v", err)
	}

	logger.InitLogger(cfg.GetLoggerDevel())

	logger.Infof("Start")

	ctx := context.Background()

	app, err := app.New(ctx, cfg)
	if err != nil {
		logger.Fatalf("app init failed: %v", err)
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	app.Run(ctx)

	logger.Infof("Exit")
}
