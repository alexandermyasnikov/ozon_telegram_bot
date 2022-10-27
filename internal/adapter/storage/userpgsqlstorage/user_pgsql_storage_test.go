package userpgsqlstorage_test

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/userpgsqlstorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

var errInternal = errors.New("internal error")

func setupSuite(ctx context.Context, tb testing.TB) (
	*userpgsqlstorage.UserPgsqlStorage, pgxmock.PgxConnIface, func(tb testing.TB),
) {
	tb.Helper()

	mock, err := pgxmock.NewConn()
	assert.NoError(tb, err)

	storage := userpgsqlstorage.New(mock)

	cls := func(tb testing.TB) {
		tb.Helper()

		mock.Close(ctx)
	}

	return storage, mock, cls
}

func TestUserPgsqlStorage_GetDefaultCurrency(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storage, mock, teardownSuite := setupSuite(ctx, t)

	defer teardownSuite(t)

	rows := pgxmock.NewRows([]string{"currency"}).
		AddRow("USD")

	mock.ExpectQuery(`SELECT currency FROM users`).
		WithArgs(int64(100)).
		WillReturnRows(rows)

	currency, err := storage.GetDefaultCurrency(ctx, entity.UserID(100))
	assert.NoError(t, err)

	assert.Equal(t, "USD", currency)
}

func TestUserPgsqlStorage_GetDefaultCurrencyError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storage, mock, teardownSuite := setupSuite(ctx, t)

	defer teardownSuite(t)

	mock.ExpectQuery(`SELECT currency FROM users`).
		WithArgs(int64(100)).
		WillReturnError(errInternal)

	_, err := storage.GetDefaultCurrency(ctx, entity.UserID(100))
	assert.Error(t, err)
}

func TestUserPgsqlStorage_UpdateDefaultCurrency(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storage, mock, teardownSuite := setupSuite(ctx, t)

	defer teardownSuite(t)

	mock.ExpectExec(`INSERT INTO users \(id, currency\)`).
		WithArgs(int64(100), "RUB").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := storage.UpdateDefaultCurrency(ctx, entity.UserID(100), "RUB")
	assert.NoError(t, err)
}

func TestUserPgsqlStorage_GetLimits(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storage, mock, teardownSuite := setupSuite(ctx, t)

	defer teardownSuite(t)

	rows := pgxmock.NewRows([]string{"day_limit", "week_limit", "month_limit"}).
		AddRow("1000", "0", "5000")

	mock.ExpectQuery(`SELECT day_limit, week_limit, month_limit FROM users`).
		WithArgs(int64(100)).
		WillReturnRows(rows)

	dayLimit, weekLimit, monthLimit, err := storage.GetLimits(ctx, entity.UserID(100))
	assert.NoError(t, err)

	assert.Equal(t, "1000", dayLimit.String())
	assert.Equal(t, "0", weekLimit.String())
	assert.Equal(t, "5000", monthLimit.String())
}

func TestUserPgsqlStorage_UpdateDayLimit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storage, mock, teardownSuite := setupSuite(ctx, t)

	defer teardownSuite(t)

	mock.ExpectExec(`INSERT INTO users \(id, day_limit\)`).
		WithArgs(int64(100), "123.45").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := storage.UpdateDayLimit(ctx, entity.UserID(100), decimal.New(12345, -2))
	assert.NoError(t, err)
}

func TestUserPgsqlStorage_UpdateWeekLimit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storage, mock, teardownSuite := setupSuite(ctx, t)

	defer teardownSuite(t)

	mock.ExpectExec(`INSERT INTO users \(id, week_limit\)`).
		WithArgs(int64(100), "123.45").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := storage.UpdateWeekLimit(ctx, entity.UserID(100), decimal.New(12345, -2))
	assert.NoError(t, err)
}

func TestUserPgsqlStorage_UpdateMonthLimit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storage, mock, teardownSuite := setupSuite(ctx, t)

	defer teardownSuite(t)

	mock.ExpectExec(`INSERT INTO users \(id, month_limit\)`).
		WithArgs(int64(100), "123.45").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := storage.UpdateMonthLimit(ctx, entity.UserID(100), decimal.New(12345, -2))
	assert.NoError(t, err)
}
