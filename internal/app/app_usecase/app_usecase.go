package appusecase

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	kafkareader "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/kafka/kafka_reader"
	kafkawriter "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/kafka/kafka_writer"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/service/ratesupdaterservicecbr"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/service/ratesupdaterserviceexchangerate"
	reportservice "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/service/report"
	currencycachestorage "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/currency_cache_storage" //nolint:lll
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/currencypgsqlstorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/expensepgsqlstorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/userpgsqlstorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/config"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	rateupdaterworker "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/worker/rate_updater_worker"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type worker interface {
	Run(context.Context)
}

type IRatesUpdaterService interface {
	Get(ctx context.Context, base string, codes []string) ([]entity.Rate, error)
}

type AppUsecase struct {
	worker         worker
	conn           *pgx.Conn
	tp             *sdktrace.TracerProvider
	metricsServer  *http.Server
	reader         *kafkareader.KafkaReader
	readerCallback kafkareader.MsgCallback
	reportClient   *reportservice.ReportClient
}

func New(ctx context.Context, cfg *config.Config) (AppUsecase, error) {
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

	reportClient := reportservice.NewReportClient(cfg.GetReportServiceAddr())

	expenseUsecase := usecase.NewExpenseUsecase(currencyStorage, userStorage, expenseStorage,
		ratesUpdaterService, reportClient, cfg)

	rateUpdaterWorker := rateupdaterworker.New(expenseUsecase, cfg)

	http.Handle("/metrics", promhttp.Handler())

	metricsServer := &http.Server{ //nolint:exhaustruct
		Addr:              cfg.GetPrometheusAddr(),
		ReadHeaderTimeout: 1 * time.Second,
	}

	facadeUsecase := usecase.New(expenseUsecase)

	reader := kafkareader.New(cfg.GetKafkaAddr(), usecase.ReadCmdState, "usecaseReader")

	writer := kafkawriter.New(cfg.GetKafkaAddr(), usecase.ProcessCmdState)

	middlewareMetricsUsecase := func(ctx context.Context, cmd *usecase.Command) error {
		startTime := time.Now()

		err = facadeUsecase.ExecuteCommand(ctx, cmd)

		duration := time.Since(startTime)

		metrics.SummaryExecuteTimeObserve(cmd.Name, duration.Seconds())
		metrics.CounterMsgInc(cmd.Name)

		return errors.Wrap(err, "middlewareMetricsUsecase")
	}

	readerCallback := func(ctx context.Context, key, value []byte) {
		var cmd usecase.Command

		err := json.Unmarshal(value, &cmd)
		if err != nil {
			logger.Errorf("can not unmarshal command: %v", err)
		}

		err = middlewareMetricsUsecase(ctx, &cmd)
		if err != nil {
			logger.Errorf("can not execute command: %v", err)
		}

		buf, err := json.Marshal(cmd)
		if err != nil {
			logger.Errorf("can not marshal command: %v", err)
		}

		err = writer.Write(ctx, []byte(cmd.Name), buf)
		if err != nil {
			logger.Errorf("can not write message: %v", err)
		}
	}

	return AppUsecase{
		worker:         rateUpdaterWorker,
		conn:           conn,
		tp:             tp,
		metricsServer:  metricsServer,
		reader:         reader,
		readerCallback: readerCallback,
		reportClient:   reportClient,
	}, nil
}

func (a *AppUsecase) Run(ctx context.Context) {
	var wg sync.WaitGroup

	go func() {
		if err := a.metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("prometheus handler start failed: %v", err)
		}
	}()

	wg.Add(1)

	go func() {
		defer wg.Done()
		a.reader.Read(ctx, a.readerCallback)
	}()

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

	a.reportClient.Close()
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
