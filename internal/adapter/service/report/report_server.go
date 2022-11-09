package reportservice

import (
	context "context"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/utils"
	"go.opentelemetry.io/otel"
)

type ExpenseStorage interface {
	Get(context.Context, entity.UserID, time.Time, time.Time) ([]entity.Expense, error)
}

type CurrencyStorage interface {
	Get(context.Context, string) (entity.Rate, error)
}

type UserStorage interface {
	GetDefaultCurrency(context.Context, entity.UserID) (string, error)
}

type Config interface {
	GetBaseCurrencyCode() string
}

type ReportServer struct {
	ReportServiceServer

	expenseStorage  ExpenseStorage
	currencyStorage CurrencyStorage
	userStorage     UserStorage
	config          Config
}

func NewReportServer(expenseStorage ExpenseStorage, currencyStorage CurrencyStorage,
	userStorage UserStorage, config Config) *ReportServer {
	return &ReportServer{ //nolint:exhaustruct
		expenseStorage:  expenseStorage,
		currencyStorage: currencyStorage,
		userStorage:     userStorage,
		config:          config,
	}
}

type ExpenseReportDTO struct {
	Category string
	Sum      decimal.Decimal
}

func (s *ReportServer) GetReport(ctx context.Context, req *Req) (*Resp, error) {
	ctx, span := otel.Tracer("ReportServer").Start(ctx, "GetReport")
	defer span.End()

	userID := entity.UserID(req.UserID)

	date, err := time.Parse(time.RFC1123, req.Date)
	if err != nil {
		return nil, errors.Wrap(err, "ExpenseUsecase.GetReport")
	}

	dateStart, dateEnd := utils.GetInterval(date, int(req.Interval))

	expenses, err := s.expenseStorage.Get(ctx, userID, dateStart, dateEnd)
	if err != nil {
		return nil, errors.Wrap(err, "ExpenseUsecase.GetReport")
	}

	currency, err := s.userStorage.GetDefaultCurrency(ctx, userID)
	if err != nil {
		currency = s.config.GetBaseCurrencyCode()
	}

	rate, err := s.currencyStorage.Get(ctx, currency)
	if err != nil {
		return nil, errors.Wrap(err, "ExpenseUsecase.GetReport")
	}

	uniq := make(map[string]int)
	expensesReport := make([]ExpenseReportDTO, 0, len(expenses))

	for _, expense := range expenses {
		ind, ok := uniq[expense.GetCategory()]

		price := expense.GetPrice().Mul(rate.GetRatio())

		if ok {
			expensesReport[ind].Sum = expensesReport[ind].Sum.Add(price)
		} else {
			uniq[expense.GetCategory()] = len(expensesReport)
			expensesReport = append(expensesReport, ExpenseReportDTO{
				Category: expense.GetCategory(),
				Sum:      price,
			})
		}
	}

	sort.Slice(expensesReport, func(i, j int) bool {
		return expensesReport[i].Category < expensesReport[j].Category
	})

	resp := &Resp{ //nolint:exhaustruct
		Currency: currency,
		Expenses: make([]*Expense, 0, len(expensesReport)),
	}

	for _, expense := range expensesReport {
		resp.Expenses = append(resp.Expenses, &Expense{ //nolint:exhaustruct
			Category: expense.Category,
			Sum:      expense.Sum.String(),
		})
	}

	return resp, nil
}
