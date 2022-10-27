package currencypgsqlstorage_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/currencypgsqlstorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

var errInternal = errors.New("internal error")

func setupSuite(ctx context.Context, tb testing.TB) (
	*currencypgsqlstorage.CurrencyPgsqlStorage, pgxmock.PgxConnIface, func(tb testing.TB),
) {
	tb.Helper()

	mock, err := pgxmock.NewConn()
	assert.NoError(tb, err)

	storage := currencypgsqlstorage.New(mock)

	cls := func(tb testing.TB) {
		tb.Helper()

		mock.Close(ctx)
	}

	return storage, mock, cls
}

func TestCurrencyPgsqlStorage_Get(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storage, mock, teardownSuite := setupSuite(ctx, t)

	defer teardownSuite(t)

	date := time.Now()

	rows := pgxmock.NewRows([]string{"code", "ratio", "time"}).
		AddRow("USD", "0.016", date)

	mock.ExpectQuery(`SELECT code, ratio, time FROM currencies`).
		WithArgs("USD").
		WillReturnRows(rows)

	rate, err := storage.Get(ctx, "USD")
	assert.NoError(t, err)

	assert.Equal(t, entity.NewRate("USD", decimal.New(16, -3), date), rate)
}

func TestCurrencyPgsqlStorage_GetUnknownCurrency(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storage, mock, teardownSuite := setupSuite(ctx, t)

	defer teardownSuite(t)

	rows := pgxmock.NewRows([]string{"code", "ratio", "time"})

	mock.ExpectQuery(`SELECT code, ratio, time FROM currencies`).
		WithArgs("USD").
		WillReturnRows(rows)

	_, err := storage.Get(ctx, "USD")
	assert.Error(t, err)
}

func TestCurrencyPgsqlStorage_GetAll(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storage, mock, teardownSuite := setupSuite(ctx, t)

	defer teardownSuite(t)

	date := time.Now()

	rows := pgxmock.NewRows([]string{"code", "ratio", "time"}).
		AddRow("RUB", "1", date).
		AddRow("USD", "0.016", date).
		AddRow("EUR", "0.017", date)

	mock.ExpectQuery(`SELECT code, ratio, time FROM currencies`).
		WillReturnRows(rows)

	rates, err := storage.GetAll(ctx)
	assert.NoError(t, err)

	assert.ElementsMatch(t,
		[]entity.Rate{
			entity.NewRate("EUR", decimal.New(17, -3), date),
			entity.NewRate("USD", decimal.New(16, -3), date),
			entity.NewRate("RUB", decimal.New(1, 0), date),
		}, rates)
}

func TestCurrencyPgsqlStorage_Update(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storage, mock, teardownSuite := setupSuite(ctx, t)

	defer teardownSuite(t)

	date := time.Now()

	mock.ExpectExec(`INSERT INTO currencies`).
		WithArgs("USD", "0.016", date).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := storage.Update(ctx, entity.NewRate("USD", decimal.New(16, -3), date))
	assert.NoError(t, err)
}

func TestCurrencyPgsqlStorage_UpdateError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storage, mock, teardownSuite := setupSuite(ctx, t)

	defer teardownSuite(t)

	date := time.Now()

	mock.ExpectExec(`INSERT INTO currencies`).
		WithArgs("USD", "0.016", date).
		WillReturnError(errInternal)

	err := storage.Update(ctx, entity.NewRate("USD", decimal.New(16, -3), date))
	assert.Error(t, err)
}
