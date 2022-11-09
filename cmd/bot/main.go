package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	appreportservice "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/app/app_report_service"
	apptgclientreader "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/app/app_tgclient_reader"
	apptgclientwriter "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/app/app_tgclient_writer"
	appusecase "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/app/app_usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/config"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
)

func main() {
	var (
		configPath = flag.String("config", config.ConfigFile, "path to config file")
		appName    = flag.String("name", "", "application name")
	)

	flag.Parse()

	cfg, err := config.New(*configPath)
	if err != nil {
		logger.Fatalf("config init failed: %v", err)
	}

	logger.InitLogger(cfg.GetLoggerDevel())

	logger.Infof("Start")

	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()

	switch *appName {
	case "tgclient_reader":
		{
			app, err := apptgclientreader.New(ctx, cfg)
			if err != nil {
				logger.Fatalf("app %v init failed: %v", *appName, err)
			}

			app.Run(ctx)
		}
	case "tgclient_writer":
		{
			app, err := apptgclientwriter.New(ctx, cfg)
			if err != nil {
				logger.Fatalf("app %v init failed: %v", *appName, err)
			}

			app.Run(ctx)
		}
	case "usecase":
		{
			app, err := appusecase.New(ctx, cfg)
			if err != nil {
				logger.Fatalf("app %v init failed: %v", *appName, err)
			}

			app.Run(ctx)
		}
	case "report_service":
		{
			app, err := appreportservice.New(ctx, cfg)
			if err != nil {
				logger.Fatalf("app %v init failed: %v", *appName, err)
			}

			app.Run(ctx)
		}
	default:
		{
			logger.Fatalf("unknown appName: %v", *appName)
		}
	}

	logger.Infof("Exit")
}
