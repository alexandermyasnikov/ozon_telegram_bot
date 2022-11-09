package reportservice

import (
	context "context"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
	"go.opentelemetry.io/otel"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReportClient struct {
	conn   *grpc.ClientConn
	client ReportServiceClient
}

func NewReportClient(addr string) *ReportClient {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("can not create grpc connection: %v", err)
	}

	client := NewReportServiceClient(conn)

	return &ReportClient{
		conn:   conn,
		client: client,
	}
}

func (c *ReportClient) Close() {
	c.conn.Close()
}

func (c *ReportClient) GetReport(ctx context.Context, req usecase.GetReportReqDTO) (usecase.GetReportRespDTO, error) {
	ctx, span := otel.Tracer("ReportClient").Start(ctx, "GetReport")
	defer span.End()

	reqRPC := &Req{ //nolint:exhaustruct
		UserID:   req.UserID,
		Date:     req.Date.Format(time.RFC1123),
		Interval: int32(req.IntervalType),
	}

	respRPC, err := c.client.GetReport(ctx, reqRPC)
	if err != nil {
		return usecase.GetReportRespDTO{}, errors.Wrap(err, "ReportClient.GetReport")
	}

	if respRPC == nil {
		logger.Errorf("GetReport: respRPC is nil")

		return usecase.GetReportRespDTO{}, errors.Wrap(err, "ReportClient.GetReport")
	}

	resp := usecase.GetReportRespDTO{
		Currency: respRPC.Currency,
		Expenses: make([]usecase.ExpenseReportDTO, 0, len(respRPC.Expenses)),
	}

	for _, expense := range respRPC.Expenses {
		sum, err := decimal.NewFromString(expense.Sum)
		if err != nil {
			return usecase.GetReportRespDTO{}, errors.Wrap(err, "ReportClient.GetReport")
		}

		resp.Expenses = append(resp.Expenses, usecase.ExpenseReportDTO{
			Category: expense.Category,
			Sum:      sum,
		})
	}

	return resp, errors.Wrap(err, "ReportClient.GetReport")
}
