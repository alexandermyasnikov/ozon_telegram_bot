package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/app"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/config"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
)

func TestAppTest(t *testing.T) { //nolint:paralleltest
	if testing.Short() {
		t.Skip("skip integration test")
	}

	cfg := &config.Config{
		Telegram: config.TelegramConfig{
			Enable: false,
			Token:  "",
		},
		Rates: config.RatesConfig{
			Service:         "cbr",
			Base:            "RUB",
			Codes:           []string{"EUR", "USD", "CNY"},
			FreqUpdateInSec: 600,
		},
		Database: config.DatabaseConfig{
			URL: "postgres://postgres:password@0.0.0.0:5432/test?sslmode=disable",
		},
		Jaeger: config.JaegerConfig{
			URL: "http://localhost:14268/api/traces",
		},
		Logger: config.LoggerConfig{
			Devel: true,
		},
		CurrencyCache: config.CacheConfig{
			Enable: true,
			Size:   10,
			TTL:    600,
		},
		ReportCache: config.CacheConfig{
			Enable: true,
			Size:   1000,
			TTL:    600,
		},
	}

	logger.InitLogger(cfg.GetLoggerDevel())

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	app, err := app.New(ctx, cfg)
	if err != nil {
		logger.Fatalf("app init failed: %v", err)
	}

	end := make(chan bool)

	go func() {
		app.Run(ctx)
		end <- true
	}()

	time.Sleep(1 * time.Second)

	type testCase struct {
		description string
		userID      int64
		textIn      string
		textOut     string
		date        time.Time
	}

	tests := [...]testCase{
		{
			description: "start",
			userID:      1,
			textIn:      `/start`,
			textOut: `Привет. Напиши свои расходы и я запомпю их.
Введи /help для более подробной информации`,
			date: time.Now(),
		},
		{
			description: "SetCurrercy",
			userID:      1,
			textIn:      `валюта USD`,
			textOut:     `Задана валюта по умолчанию USD`,
			date:        time.Now(),
		},
		{
			description: "GetLimits",
			userID:      1,
			textIn:      `лимиты`,
			textOut: `Текущие лимиты:
Дневной - 0.00 USD
Недельный - 0.00 USD
Месячный - 0.00 USD`,
			date: time.Now(),
		},
		{
			description: "SetLimit",
			userID:      1,
			textIn:      `лимит неделя 50`,
			textOut:     `Установил лимит: неделя - 50.00 - USD`,
			date:        time.Now(),
		},
		{
			description: "GetLimits",
			userID:      1,
			textIn:      `лимиты`,
			textOut: `Текущие лимиты:
Дневной - 0.00 USD
Недельный - 50.00 USD
Месячный - 0.00 USD`,
			date: time.Now(),
		},
		{
			description: "AddExpense",
			userID:      1,
			textIn:      `расход AppStore 4.50`,
			textOut:     `Добавил AppStore - 4.50 USD Sat, 29 Oct 2022 19:34:28 UTC`,
			date:        time.Date(2022, 10, 29, 19, 34, 28, 0, time.UTC),
		},
		{
			description: "AddExpense",
			userID:      1,
			textIn:      `расход AppleTV 46.00`,
			textOut: `Добавил AppleTV - 46.00 USD Sat, 29 Oct 2022 19:34:30 UTC
Внимание! Превышен лимит: неделя - 0.50`,
			date: time.Date(2022, 10, 29, 19, 34, 30, 0, time.UTC),
		},
		{
			description: "GetReport",
			userID:      1,
			textIn:      `отчет день`,
			textOut: `Расходы по категориям за день:
AppStore - 4.50
AppleTV - 46.00`,
			date: time.Date(2022, 10, 29, 19, 34, 32, 0, time.UTC),
		},
	}

	for _, scenario := range tests { //nolint:paralleltest
		t.Run(scenario.description, func(t *testing.T) {
			textActual, err := app.GetRouterForTest().Execute(ctx, scenario.userID, scenario.textIn, scenario.date)
			assert.NoError(t, err)
			assert.Equal(t, scenario.textOut, textActual)
		})
	}

	cancel()

	<-end
}
