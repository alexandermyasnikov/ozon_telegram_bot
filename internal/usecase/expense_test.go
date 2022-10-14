package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/usecase/mock_usecase"
)

var errUnknown = errors.New("unknown error")

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
		userStorage.EXPECT().UpdateDefaultCurrency(entity.NewUserID(201), "RUB").Return(nil),
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
		userStorage.EXPECT().UpdateDefaultCurrency(entity.NewUserID(201), "USD").Return(nil),
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
		userStorage.EXPECT().UpdateDefaultCurrency(entity.NewUserID(201), "RUB").Return(errUnknown),
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
				entity.NewRate("RUB", entity.NewDecimal(1, 0), entity.NewDateTime(2022, 10, 1, 0, 0, 0)),
				entity.NewRate("EUR", entity.NewDecimal(16, 3), entity.NewDateTime(2022, 10, 1, 0, 0, 0)),
			}, nil),
		currencyStorage.EXPECT().Update(
			entity.NewRate("RUB", entity.NewDecimal(1, 0), entity.NewDateTime(2022, 10, 1, 0, 0, 0)),
		).Return(nil),
		currencyStorage.EXPECT().Update(
			entity.NewRate("EUR", entity.NewDecimal(16, 3), entity.NewDateTime(2022, 10, 1, 0, 0, 0)),
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

	time1 := time.Date(2022, time.October, 10, 0, 0, 0, 0, &time.Location{})

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	currencyStorage := mock_usecase.NewMockICurrencyStorage(ctrl)
	userStorage := mock_usecase.NewMockIUserStorage(ctrl)
	expenseStorage := mock_usecase.NewMockIExpenseStorage(ctrl)
	ratesUpdaterService := mock_usecase.NewMockIRatesUpdaterService(ctrl)
	config := mock_usecase.NewMockIConfig(ctrl)
	expenseUsecase := usecase.NewExpenseUsecase(currencyStorage, userStorage, expenseStorage, ratesUpdaterService, config)

	config.EXPECT().GetBaseCurrencyCode().Return("RUB").AnyTimes()
	config.EXPECT().GetCurrencyCodes().Return([]string{"USD", "EUR", "JPY"}).AnyTimes()
	config.EXPECT().GetFrequencyRateUpdateSec().Return(600).AnyTimes()

	currencyStorage.EXPECT().Get("RUB").Return(
		entity.NewRate("RUB", entity.NewDecimal(1, 0), entity.NewDateTimeFromTime(time.Now())),
		nil)
	currencyStorage.EXPECT().Get("EUR").Return(
		entity.NewRate("EUR", entity.NewDecimal(16, 3), entity.NewDateTimeFromTime(time.Now())), nil)
	expenseStorage.EXPECT().Create(entity.NewUserID(202),
		entity.NewExpense("Netflix", entity.NewDecimal(625, 0), entity.NewDateFromTime(time1))).Return(nil)

	req := usecase.AddExpenseReqDTO{
		UserID:   202,
		Category: "Netflix",
		Price:    10,
		Date:     time1,
		Currency: "EUR",
	}

	err := expenseUsecase.AddExpense(ctx, req)
	assert.NoError(t, err)
}

func TestGetReport(t *testing.T) {
	t.Parallel()

	date := entity.NewDate(2022, 10, 1)

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	currencyStorage := mock_usecase.NewMockICurrencyStorage(ctrl)
	userStorage := mock_usecase.NewMockIUserStorage(ctrl)
	expenseStorage := mock_usecase.NewMockIExpenseStorage(ctrl)
	ratesUpdaterService := mock_usecase.NewMockIRatesUpdaterService(ctrl)
	config := mock_usecase.NewMockIConfig(ctrl)
	expenseUsecase := usecase.NewExpenseUsecase(currencyStorage, userStorage, expenseStorage, ratesUpdaterService, config)

	config.EXPECT().GetBaseCurrencyCode().Return("RUB").AnyTimes()
	config.EXPECT().GetCurrencyCodes().Return([]string{"USD", "EUR", "JPY"}).AnyTimes()
	config.EXPECT().GetFrequencyRateUpdateSec().Return(600).AnyTimes()

	currencyStorage.EXPECT().Get("RUB").Return(
		entity.NewRate("RUB", entity.NewDecimal(1, 0), entity.NewDateTimeFromTime(time.Now())), nil)

	expenseStorage.EXPECT().Get(entity.NewUserID(202), date, 7).Return(
		[]entity.Expense{
			entity.NewExpense("Spotify", entity.NewDecimal(600, 0), date),
			entity.NewExpense("appStore", entity.NewDecimal(601345, 3), date),
			entity.NewExpense("appStore", entity.NewDecimal(1059, 2), date),
		}, nil)
	currencyStorage.EXPECT().Get("EUR").Return(
		entity.NewRate("EUR", entity.NewDecimal(16, 3), entity.NewDateTime(2022, 10, 1, 0, 0, 0)), nil)

	req := usecase.GetReportReqDTO{
		UserID:   202,
		Date:     date.ToTime(),
		Days:     7,
		Currency: "EUR",
	}

	resp, err := expenseUsecase.GetReport(ctx, req)
	assert.NoError(t, err)

	assert.EqualValues(t, usecase.GetReportRespDTO{
		Currency: "EUR",
		Categories: map[string]entity.Decimal{
			"Spotify":  entity.NewDecimal(96, 1),
			"appStore": entity.NewDecimal(979096, 5),
		}}, resp)
}
