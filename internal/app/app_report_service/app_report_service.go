package appreportservice

import (
	"context"
	"net"

	"github.com/jackc/pgx/v5"
	reportservice "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/service/report"
	currencycachestorage "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/currency_cache_storage" //nolint:lll
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/currencypgsqlstorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/expensepgsqlstorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/userpgsqlstorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/config"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
	"google.golang.org/grpc"
)

type AppUsecase struct {
	conn     *pgx.Conn
	listener net.Listener
	server   *grpc.Server
}

func New(ctx context.Context, cfg *config.Config) (AppUsecase, error) {
	conn, err := pgx.Connect(ctx, cfg.GetDatabaseURL())
	if err != nil {
		logger.Fatalf("pg client init failed: %v", err)
	}

	listener, err := net.Listen("tcp", cfg.GetReportServiceAddr())
	if err != nil {
		logger.Fatalf("can not create listener: %v", err)
	}

	currencyStorage := currencycachestorage.New(currencypgsqlstorage.New(conn), cfg)
	userStorage := userpgsqlstorage.New(conn)
	expenseStorage := expensepgsqlstorage.New(conn)

	reportService := reportservice.NewReportServer(expenseStorage, currencyStorage, userStorage, cfg)

	server := grpc.NewServer()
	reportservice.RegisterReportServiceServer(server, reportService)

	return AppUsecase{
		conn:     conn,
		listener: listener,
		server:   server,
	}, nil
}

func (a *AppUsecase) Run(ctx context.Context) {
	go func() {
		if err := a.server.Serve(a.listener); err != nil {
			logger.Fatalf("can not run server: %v", err)
		}
	}()

	<-ctx.Done()

	a.server.Stop()

	a.conn.Close(ctx)
}
