package app

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/service/ratesupdaterservicecbr"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/service/ratesupdaterserviceexchangerate"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/currencypgsqlstorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/expensepgsqlstorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/userpgsqlstorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter/texthandler"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	rateupdaterworker "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/worker/rate_updater_worker"
)

type client interface {
	Run(context.Context)
}

type worker interface {
	Run(context.Context)
}

type IRatesUpdaterService interface {
	Get(ctx context.Context, base string, codes []string) ([]entity.Rate, error)
}

type config interface {
	TelegramToken() string
	GetBaseCurrencyCode() string
	GetCurrencyCodes() []string
	GetFrequencyRateUpdateSec() int
	GetRatesService() string
	GetDatabaseURL() string
}

var _ client = &tg.Client{}

type App struct {
	client client
	worker worker
	conn   *pgx.Conn
}

func New(cfg config) (App, error) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.GetDatabaseURL())
	if err != nil {
		log.Fatal("pg client init failed:", err)
	}

	currencyStorage := currencypgsqlstorage.New(conn)
	userStorage := userpgsqlstorage.New(conn)
	expenseStorage := expensepgsqlstorage.New(conn)

	var ratesUpdaterService IRatesUpdaterService

	if cfg.GetRatesService() == "cbr" {
		ratesUpdaterService = ratesupdaterservicecbr.New()
	} else {
		ratesUpdaterService = ratesupdaterserviceexchangerate.New()
	}

	expenseUsecase := usecase.NewExpenseUsecase(currencyStorage, userStorage, expenseStorage,
		ratesUpdaterService, cfg)

	routerText := textrouter.New()

	routerText.Register(texthandler.NewStart())
	routerText.Register(texthandler.NewHelp())
	routerText.Register(texthandler.NewAbout())
	routerText.Register(texthandler.NewSetDefaultCurrency(expenseUsecase))
	routerText.Register(texthandler.NewAddExpense(expenseUsecase))
	routerText.Register(texthandler.NewGetReport(expenseUsecase))
	routerText.Register(texthandler.NewSetLimit(expenseUsecase))
	routerText.Register(texthandler.NewGetLimits(expenseUsecase))
	routerText.Register(texthandler.NewUnknown())

	rateUpdaterWorker := rateupdaterworker.New(expenseUsecase, cfg)

	tgClient, err := tg.New(cfg, routerText)
	if err != nil {
		log.Fatal("tg client init failed:", err)
	}

	return App{
		client: tgClient,
		worker: rateUpdaterWorker,
		conn:   conn,
	}, nil
}

func (a *App) Run(ctx context.Context) {
	var wg sync.WaitGroup

	wg.Add(1 + 1)

	go func() {
		defer wg.Done()
		a.client.Run(ctx)
	}()

	go func() {
		defer wg.Done()
		a.worker.Run(ctx)
	}()

	wg.Wait()

	a.conn.Close(ctx)
}
