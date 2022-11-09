package appusecase_test

import (
	"context"
	"database/sql"
	"sync"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/assert"
	appreportservice "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/app/app_report_service"
	apptgclientreader "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/app/app_tgclient_reader"
	apptgclientwriter "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/app/app_tgclient_writer"
	appusecase "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/app/app_usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/config"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
	fakeclientreader "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/test/client/fake_client_reader"
	fakeclientwriter "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/test/client/fake_client_writer"
)

// Запускаем все сервисы, кроме телеграм клиента
// Выполняем тестовый сценарий

func TestApp_Integration(t *testing.T) { //nolint:paralleltest,maintidx
	if testing.Short() {
		t.Skip("skip integration test")
	}

	cfg := &config.Config{
		Telegram: config.TelegramConfig{
			Token: "",
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
		Kafka: config.KafkaConfig{
			Addr: "0.0.0.0:9092",
		},
		Prometheus: config.PrometheusConfig{
			Addr: "0.0.0.0:8080",
		},
		ReportService: config.ReportServiceConfig{
			Addr: "0.0.0.0:9094",
		},
	}

	logger.InitLogger(cfg.GetLoggerDevel())

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	conn, err := sql.Open("pgx", cfg.GetDatabaseURL())
	assert.NoError(t, err)
	err = goose.Reset(conn, "../db/migrations")
	assert.NoError(t, err)
	err = goose.Up(conn, "../db/migrations")
	assert.NoError(t, err)
	conn.Close()

	timeHelper := func(offset int) time.Time {
		date, _ := time.Parse(time.RFC1123, "Wed, 09 Nov 2022 16:00:00 MSK")

		return date.Add(time.Duration(offset) * time.Minute)
	}

	type testCase struct {
		description  string
		userID       int64
		date         time.Time
		text         string
		textExpected string
	}

	tests := [...]testCase{
		{
			description: "start",
			userID:      1,
			date:        timeHelper(0),
			text:        `/start`,
			textExpected: `Привет. Напиши свои расходы и я запомпю их.
Введи /help для более подробной информации`,
		},

		{
			description:  "setCurrency",
			userID:       1,
			date:         timeHelper(10),
			text:         `валюта USD`,
			textExpected: `Задана валюта по умолчанию USD`,
		},
		{
			description: "getLimits",
			userID:      1,
			date:        timeHelper(11),
			text:        `лимиты`,
			textExpected: `Текущие лимиты:
Дневной - 0.00 USD
Недельный - 0.00 USD
Месячный - 0.00 USD`,
		},

		{
			description:  "setLimit",
			userID:       1,
			date:         timeHelper(20),
			text:         `лимит день 10`,
			textExpected: `Установил лимит: день - 10.00 - USD`,
		},
		{
			description:  "setLimit",
			userID:       1,
			date:         timeHelper(21),
			text:         `лимит неделя 50`,
			textExpected: `Установил лимит: неделя - 50.00 - USD`,
		},
		{
			description: "GetLimits",
			userID:      1,
			date:        timeHelper(22),
			text:        `лимиты`,
			textExpected: `Текущие лимиты:
Дневной - 10.00 USD
Недельный - 50.00 USD
Месячный - 0.00 USD`,
		},

		{
			description:  "AddExpense",
			userID:       1,
			date:         timeHelper(30),
			text:         `расход Netflix 4.50`,
			textExpected: `Добавил Netflix - 4.50 USD Wed, 09 Nov 2022 16:30:00 MSK`,
		},
		{
			description: "GetReportDay",
			userID:      1,
			date:        timeHelper(31),
			text:        `отчет день`,
			textExpected: `Расходы по категориям за день:
Netflix - 4.50`,
		},
		{
			description: "GetReporWeek",
			userID:      1,
			date:        timeHelper(32),
			text:        `отчет неделя`,
			textExpected: `Расходы по категориям за неделя:
Netflix - 4.50`,
		},

		{
			description:  "AddExpense",
			userID:       1,
			date:         timeHelper(40),
			text:         `расход AppStore 5.00`,
			textExpected: `Добавил AppStore - 5.00 USD Wed, 09 Nov 2022 16:40:00 MSK`,
		},
		{
			description: "GetReportDay",
			userID:      1,
			date:        timeHelper(41),
			text:        `отчет день`,
			textExpected: `Расходы по категориям за день:
AppStore - 5.00
Netflix - 4.50`,
		},
		{
			description: "GetReporWeek",
			userID:      1,
			date:        timeHelper(42),
			text:        `отчет неделя`,
			textExpected: `Расходы по категориям за неделя:
AppStore - 5.00
Netflix - 4.50`,
		},

		{
			description: "AddExpense",
			userID:      1,
			date:        timeHelper(50),
			text:        `расход AppStore 2.00`,
			textExpected: `Добавил AppStore - 2.00 USD Wed, 09 Nov 2022 16:50:00 MSK
Внимание! Превышен лимит: день - 1.50`,
		},
		{
			description: "GetReportDay",
			userID:      1,
			date:        timeHelper(51),
			text:        `отчет день`,
			textExpected: `Расходы по категориям за день:
AppStore - 7.00
Netflix - 4.50`,
		},
		{
			description: "GetReporWeek",
			userID:      1,
			date:        timeHelper(52),
			text:        `отчет неделя`,
			textExpected: `Расходы по категориям за неделя:
AppStore - 7.00
Netflix - 4.50`,
		},

		{
			description:  "AddExpense",
			userID:       1,
			date:         timeHelper(24*60 + 60),
			text:         `расход Food 6.00`,
			textExpected: `Добавил Food - 6.00 USD Thu, 10 Nov 2022 17:00:00 MSK`,
		},
		{
			description: "GetReportDay",
			userID:      1,
			date:        timeHelper(24*60 + 61),
			text:        `отчет день`,
			textExpected: `Расходы по категориям за день:
Food - 6.00`,
		},
		{
			description: "GetReporWeek",
			userID:      1,
			date:        timeHelper(24*60 + 62),
			text:        `отчет неделя`,
			textExpected: `Расходы по категориям за неделя:
AppStore - 7.00
Food - 6.00
Netflix - 4.50`,
		},

		{
			description: "AddExpense",
			userID:      1,
			date:        timeHelper(24*60 + 70),
			text:        `расход Steam 100.00`,
			textExpected: `Добавил Steam - 100.00 USD Thu, 10 Nov 2022 17:10:00 MSK
Внимание! Превышен лимит: день - 96.00
Внимание! Превышен лимит: неделя - 67.50`,
		},
		{
			description: "GetReportDay",
			userID:      1,
			date:        timeHelper(24*60 + 71),
			text:        `отчет день`,
			textExpected: `Расходы по категориям за день:
Food - 6.00
Steam - 100.00`,
		},
		{
			description: "GetReporWeek",
			userID:      1,
			date:        timeHelper(24*60 + 72),
			text:        `отчет неделя`,
			textExpected: `Расходы по категориям за неделя:
AppStore - 7.00
Food - 6.00
Netflix - 4.50
Steam - 100.00`,
		},
	}

	messages := make([]fakeclientreader.Message, 0)
	for _, scenario := range tests {
		messages = append(messages, fakeclientreader.Message{
			UserID: scenario.userID,
			Date:   scenario.date,
			Text:   scenario.text,
		})
	}

	appClientReader, err := apptgclientreader.NewWithCustomClient(ctx, cfg, fakeclientreader.New(messages, 1*time.Second))
	if err != nil {
		logger.Fatalf("app tgclient init failed: %v", err)
	}

	appUsecase, err := appusecase.New(ctx, cfg)
	if err != nil {
		logger.Fatalf("app usecase init failed: %v", err)
	}

	appReportService, err := appreportservice.New(ctx, cfg)
	if err != nil {
		logger.Fatalf("app reportService init failed: %v", err)
	}

	clientWriter := fakeclientwriter.New()

	appClientWriter, err := apptgclientwriter.NewWithCustomClient(ctx, cfg, clientWriter)
	if err != nil {
		logger.Fatalf("app tgclient init failed: %v", err)
	}

	func() {
		var wg sync.WaitGroup

		wg.Add(1 + 1 + 1 + 1)

		go func() {
			defer wg.Done()
			appClientReader.Run(ctx)
			time.Sleep(10 * time.Second)
			cancel()
		}()

		go func() {
			defer wg.Done()
			appUsecase.Run(ctx)
		}()

		go func() {
			defer wg.Done()
			appReportService.Run(ctx)
		}()

		go func() {
			defer wg.Done()
			appClientWriter.Run(ctx)
		}()

		wg.Wait()
	}()

	messagesAct := clientWriter.GetMessages()

	assert.Equal(t, len(tests), len(messagesAct))

	for i, scenario := range tests { //nolint:paralleltest
		t.Run(scenario.description, func(t *testing.T) {
			assert.Equal(t, scenario.userID, messagesAct[i].UserID)
			assert.Equal(t, scenario.textExpected, messagesAct[i].Text)
		})
	}
}
