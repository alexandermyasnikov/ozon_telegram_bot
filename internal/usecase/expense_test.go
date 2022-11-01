package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase/mock_usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/utils"
)

var errUnknown = errors.New("unknown error")

var timeHelper = func(year, month, day int) time.Time { //nolint:gochecknoglobals
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func TestExpenseSetDefaultCurrency_CurrencyEqBaseCode(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	currencyStorage := mock_usecase.NewMockICurrencyStorage(ctrl)
	userStorage := mock_usecase.NewMockIUserStorage(ctrl)
	expenseStorage := mock_usecase.NewMockIExpenseStorage(ctrl)
	ratesUpdaterService := mock_usecase.NewMockIRatesUpdaterService(ctrl)
	config := mock_usecase.NewMockIConfig(ctrl)
	expenseUsecase := usecase.NewExpenseUsecase(currencyStorage, userStorage, expenseStorage, ratesUpdaterService, config)

	gomock.InOrder(
		config.EXPECT().GetBaseCurrencyCode().Return("RUB"),
		userStorage.EXPECT().UpdateDefaultCurrency(gomock.Any(), entity.UserID(201), "RUB").Return(nil),
	)

	req := usecase.SetDefaultCurrencyReqDTO{
		UserID:   201,
		Currency: "RUB",
	}

	err := expenseUsecase.SetDefaultCurrency(ctx, req)
	assert.NoError(t, err)
}

func TestExpenseSetDefaultCurrency_CurrencyInCodes(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	currencyStorage := mock_usecase.NewMockICurrencyStorage(ctrl)
	userStorage := mock_usecase.NewMockIUserStorage(ctrl)
	expenseStorage := mock_usecase.NewMockIExpenseStorage(ctrl)
	ratesUpdaterService := mock_usecase.NewMockIRatesUpdaterService(ctrl)
	config := mock_usecase.NewMockIConfig(ctrl)
	expenseUsecase := usecase.NewExpenseUsecase(currencyStorage, userStorage, expenseStorage, ratesUpdaterService, config)

	gomock.InOrder(
		config.EXPECT().GetBaseCurrencyCode().Return("RUB"),
		config.EXPECT().GetCurrencyCodes().Return([]string{"CNY", "EUR", "USD", "JPY"}),
		userStorage.EXPECT().UpdateDefaultCurrency(gomock.Any(), entity.UserID(201), "USD").Return(nil),
	)

	req := usecase.SetDefaultCurrencyReqDTO{
		UserID:   201,
		Currency: "USD",
	}

	err := expenseUsecase.SetDefaultCurrency(ctx, req)
	assert.NoError(t, err)
}

func TestExpenseSetDefaultCurrency_UnknownCurrency(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	currencyStorage := mock_usecase.NewMockICurrencyStorage(ctrl)
	userStorage := mock_usecase.NewMockIUserStorage(ctrl)
	expenseStorage := mock_usecase.NewMockIExpenseStorage(ctrl)
	ratesUpdaterService := mock_usecase.NewMockIRatesUpdaterService(ctrl)
	config := mock_usecase.NewMockIConfig(ctrl)
	expenseUsecase := usecase.NewExpenseUsecase(currencyStorage, userStorage, expenseStorage, ratesUpdaterService, config)

	gomock.InOrder(
		config.EXPECT().GetBaseCurrencyCode().Return("RUB"),
		config.EXPECT().GetCurrencyCodes().Return([]string{"CNY", "EUR", "USD", "JPY"}),
	)

	req := usecase.SetDefaultCurrencyReqDTO{
		UserID:   201,
		Currency: "KZT",
	}

	err := expenseUsecase.SetDefaultCurrency(ctx, req)
	assert.Error(t, err)
	assert.EqualError(t, err, "currency is unsupported")
}

func TestExpenseSetDefaultCurrency_DBError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	currencyStorage := mock_usecase.NewMockICurrencyStorage(ctrl)
	userStorage := mock_usecase.NewMockIUserStorage(ctrl)
	expenseStorage := mock_usecase.NewMockIExpenseStorage(ctrl)
	ratesUpdaterService := mock_usecase.NewMockIRatesUpdaterService(ctrl)
	config := mock_usecase.NewMockIConfig(ctrl)
	expenseUsecase := usecase.NewExpenseUsecase(currencyStorage, userStorage, expenseStorage, ratesUpdaterService, config)

	gomock.InOrder(
		config.EXPECT().GetBaseCurrencyCode().Return("RUB"),
		userStorage.EXPECT().UpdateDefaultCurrency(gomock.Any(), entity.UserID(201), "RUB").Return(errUnknown),
	)

	req := usecase.SetDefaultCurrencyReqDTO{
		UserID:   201,
		Currency: "RUB",
	}

	err := expenseUsecase.SetDefaultCurrency(ctx, req)
	assert.Error(t, err)
	assert.EqualError(t, err, "ExpenseUsecase.SetDefaultCurrency: unknown error")
}

func TestUpdateCurrency(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	currencyStorage := mock_usecase.NewMockICurrencyStorage(ctrl)
	userStorage := mock_usecase.NewMockIUserStorage(ctrl)
	expenseStorage := mock_usecase.NewMockIExpenseStorage(ctrl)
	ratesUpdaterService := mock_usecase.NewMockIRatesUpdaterService(ctrl)
	config := mock_usecase.NewMockIConfig(ctrl)
	expenseUsecase := usecase.NewExpenseUsecase(currencyStorage, userStorage, expenseStorage, ratesUpdaterService, config)

	gomock.InOrder(
		config.EXPECT().GetBaseCurrencyCode().Return("RUB"),
		config.EXPECT().GetCurrencyCodes().Return([]string{"USD", "EUR"}),
		ratesUpdaterService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(
			[]entity.Rate{
				entity.NewRate("RUB", decimal.New(1, 0), timeHelper(2022, 10, 1)),
				entity.NewRate("EUR", decimal.New(16, 3), timeHelper(2022, 10, 1)),
			}, nil),
		currencyStorage.EXPECT().Update(gomock.Any(),
			entity.NewRate("RUB", decimal.New(1, 0), timeHelper(2022, 10, 1)),
		).Return(nil),
		currencyStorage.EXPECT().Update(gomock.Any(),
			entity.NewRate("EUR", decimal.New(16, 3), timeHelper(2022, 10, 1)),
		).Return(nil),
	)

	err := expenseUsecase.UpdateCurrency(ctx)
	assert.NoError(t, err)
}

func TestUpdateCurrency_SrvError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	currencyStorage := mock_usecase.NewMockICurrencyStorage(ctrl)
	userStorage := mock_usecase.NewMockIUserStorage(ctrl)
	expenseStorage := mock_usecase.NewMockIExpenseStorage(ctrl)
	ratesUpdaterService := mock_usecase.NewMockIRatesUpdaterService(ctrl)
	config := mock_usecase.NewMockIConfig(ctrl)
	expenseUsecase := usecase.NewExpenseUsecase(currencyStorage, userStorage, expenseStorage, ratesUpdaterService, config)

	gomock.InOrder(
		config.EXPECT().GetBaseCurrencyCode().Return("RUB"),
		config.EXPECT().GetCurrencyCodes().Return([]string{"USD", "EUR"}),
		ratesUpdaterService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errUnknown),
	)

	err := expenseUsecase.UpdateCurrency(ctx)
	assert.Error(t, err)
	assert.EqualError(t, err, "ExpenseUsecase.UpdateCurrency: ExpenseUsecase.tryUpdateRates: unknown error")
}

func TestAddExpense(t *testing.T) {
	t.Parallel()

	time1 := timeHelper(2022, 11, 10)

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	currencyStorage := mock_usecase.NewMockICurrencyStorage(ctrl)
	userStorage := mock_usecase.NewMockIUserStorage(ctrl)
	expenseStorage := mock_usecase.NewMockIExpenseStorage(ctrl)
	ratesUpdaterService := mock_usecase.NewMockIRatesUpdaterService(ctrl)
	config := mock_usecase.NewMockIConfig(ctrl)
	expenseUsecase := usecase.NewExpenseUsecase(currencyStorage, userStorage, expenseStorage, ratesUpdaterService, config)

	gomock.InOrder(
		config.EXPECT().GetBaseCurrencyCode().Return("RUB"),
		currencyStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(entity.Rate{}, nil),
		config.EXPECT().GetFrequencyRateUpdateSec().Return(60),
		config.EXPECT().GetBaseCurrencyCode().Return("RUB"),
		config.EXPECT().GetCurrencyCodes().Return([]string{"USD", "EUR", "JPY"}),
		ratesUpdaterService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil),

		userStorage.EXPECT().GetDefaultCurrency(gomock.Any(), gomock.Any()).
			Return("EUR", nil),
		currencyStorage.EXPECT().Get(gomock.Any(), "EUR").
			Return(entity.NewRate("EUR", decimal.New(16, -3), time.Now()), nil),
		expenseStorage.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil),
		userStorage.EXPECT().GetLimits(gomock.Any(), gomock.Any()).
			Return(decimal.New(10, 0), decimal.New(0, 0), decimal.New(50, 0), nil),
		expenseStorage.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, nil),
		expenseStorage.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return([]entity.Expense{
				entity.NewExpense("Category1", decimal.New(2, 0), time1),
			}, nil),
	)

	req := usecase.AddExpenseReqDTO{
		UserID:   202,
		Category: "Netflix",
		Price:    decimal.New(10, 0),
		Date:     time1,
	}

	resp, err := expenseUsecase.AddExpense(ctx, req)
	assert.NoError(t, err)

	assert.EqualValues(t, usecase.AddExpenseRespDTO{
		Currency: "EUR",
		Limits: map[int]decimal.Decimal{
			utils.DayInterval:   decimal.New(160, -3),
			utils.MonthInterval: decimal.New(768, -3),
		},
	}, resp)
}

func TestGetReport(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	currencyStorage := mock_usecase.NewMockICurrencyStorage(ctrl)
	userStorage := mock_usecase.NewMockIUserStorage(ctrl)
	expenseStorage := mock_usecase.NewMockIExpenseStorage(ctrl)
	ratesUpdaterService := mock_usecase.NewMockIRatesUpdaterService(ctrl)
	config := mock_usecase.NewMockIConfig(ctrl)
	expenseUsecase := usecase.NewExpenseUsecase(currencyStorage, userStorage, expenseStorage, ratesUpdaterService, config)

	gomock.InOrder(
		config.EXPECT().GetBaseCurrencyCode().Return("RUB"),
		currencyStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(entity.Rate{}, nil),
		config.EXPECT().GetFrequencyRateUpdateSec().Return(60),
		config.EXPECT().GetBaseCurrencyCode().Return("RUB"),
		config.EXPECT().GetCurrencyCodes().Return([]string{"USD", "EUR", "JPY"}),
		ratesUpdaterService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil),

		expenseStorage.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return([]entity.Expense{
				entity.NewExpense("AppStore", decimal.New(3125, -1), timeHelper(2022, 11, 1)),
				entity.NewExpense("Spotify", decimal.New(125, 0), timeHelper(2022, 11, 1)),
				entity.NewExpense("AppStore", decimal.New(625, -1), timeHelper(2022, 11, 2)),
			}, nil),
		userStorage.EXPECT().GetDefaultCurrency(gomock.Any(), gomock.Any()).
			Return("EUR", nil),
		currencyStorage.EXPECT().Get(gomock.Any(), "EUR").
			Return(entity.NewRate("EUR", decimal.New(16, -3), time.Now()), nil),
	)

	req := usecase.GetReportReqDTO{
		UserID:       202,
		Date:         timeHelper(2022, 10, 1),
		IntervalType: utils.WeekInterval,
	}

	resp, err := expenseUsecase.GetReport(ctx, req)
	assert.NoError(t, err)

	assert.EqualValues(t, usecase.GetReportRespDTO{
		Currency: "EUR",
		Expenses: []usecase.ExpenseReportDTO{
			{
				Category: "AppStore",
				Sum:      decimal.New(60000, -4),
			},
			{
				Category: "Spotify",
				Sum:      decimal.New(2000, -3),
			},
		},
	}, resp)
}
