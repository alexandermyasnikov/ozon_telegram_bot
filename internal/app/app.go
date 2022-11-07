package app

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/service/ratesupdaterservicecbr"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/service/ratesupdaterserviceexchangerate"
	currencycachestorage "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/currency_cache_storage" //nolint:lll
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/currencypgsqlstorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/expensepgsqlstorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/userpgsqlstorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/clients/tg"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/config"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/textrouter/texthandler"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	rateupdaterworker "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/worker/rate_updater_worker"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
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

var _ client = &tg.Client{}

type App struct {
	client        client
	worker        worker
	conn          *pgx.Conn
	router        *textrouter.RouterText
	tp            *sdktrace.TracerProvider
	metricsServer *http.Server
}

func New(ctx context.Context, cfg *config.Config) (App, error) {
	tp, err := registerJaeger(cfg.GetJaegerURL())
	if err != nil {
		logger.Fatalf("jaeger client init failed: %v", err)
	}

	conn, err := pgx.Connect(ctx, cfg.GetDatabaseURL())
	if err != nil {
		logger.Fatalf("pg client init failed: %v", err)
	}

	currencyStorage := currencycachestorage.New(currencypgsqlstorage.New(conn), cfg)

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

	var tgClient *tg.Client

	if cfg.TelegramEnable() {
		tgClient, err = tg.New(cfg, routerText)
		if err != nil {
			logger.Fatalf("tg client init failed: %v", err)
		}
	}

	http.Handle("/metrics", promhttp.Handler())

	metricsServer := &http.Server{ //nolint:exhaustruct
		Addr:              ":8080",
		ReadHeaderTimeout: 1 * time.Second,
	}

	return App{
		client:        tgClient,
		worker:        rateUpdaterWorker,
		conn:          conn,
		router:        routerText,
		tp:            tp,
		metricsServer: metricsServer,
	}, nil
}

func (a *App) Run(ctx context.Context) {
	var wg sync.WaitGroup

	go func() {
		if err := a.metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("prometheus handler start failed: %v", err)
		}
	}()

	if p, ok := a.client.(*tg.Client); ok && p != nil {
		wg.Add(1)

		go func() {
			defer wg.Done()
			a.client.Run(ctx)
		}()
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		a.worker.Run(ctx)
	}()

	wg.Wait()

	a.conn.Close(ctx)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second) //nolint:contextcheck
	defer cancel()

	if err := a.tp.Shutdown(ctx); err != nil {
		logger.Errorf("can not shutdown jaeger: %v", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second) //nolint:contextcheck
	defer cancel()

	if err := a.metricsServer.Shutdown(ctx); err != nil {
		logger.Errorf("can not shutdown metricsServer: %v", err)
	}
}

func (a *App) GetRouterForTest() *textrouter.RouterText {
	return a.router
}

func registerJaeger(url string) (*sdktrace.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, errors.Wrap(err, "RegisterJaeger")
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("telegram-bot"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{}))

	return tp, nil
}
