package expensepgsqlstorage_test

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/expensepgsqlstorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

var errInternal = errors.New("internal error")

func setupSuite(ctx context.Context, tb testing.TB) (
	*expensepgsqlstorage.ExpensePgsqlStorage, pgxmock.PgxConnIface, func(tb testing.TB),
) {
	tb.Helper()

	mock, err := pgxmock.NewConn()
	assert.NoError(tb, err)

	storage := expensepgsqlstorage.New(mock)

	cls := func(tb testing.TB) {
		tb.Helper()

		mock.Close(ctx)
	}

	return storage, mock, cls
}

func TestExpensePgsqlStorage_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storage, mock, teardownSuite := setupSuite(ctx, t)

	defer teardownSuite(t)

	date := time.Now()

	mock.ExpectExec(`INSERT INTO epxenses \(user_id, category, price, time\)`).
		WithArgs(int64(100), "Macbook", "150350.56", date).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := storage.Create(ctx, entity.UserID(100), entity.NewExpense("Macbook", decimal.New(15035056, -2), date))
	assert.NoError(t, err)
}

func TestExpensePgsqlStorage_CreateError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storage, mock, teardownSuite := setupSuite(ctx, t)

	defer teardownSuite(t)

	date := time.Now()

	mock.ExpectExec(`INSERT INTO epxenses \(user_id, category, price, time\)`).
		WithArgs(int64(100), "Macbook", "150350.56", date).
		WillReturnError(errInternal)

	err := storage.Create(ctx, entity.UserID(100), entity.NewExpense("Macbook", decimal.New(15035056, -2), date))
	assert.Error(t, err)
}

func TestExpensePgsqlStorage_GetLimits(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	storage, mock, teardownSuite := setupSuite(ctx, t)

	defer teardownSuite(t)

	dateStart := time.Now()
	dateEnd := dateStart.AddDate(0, 1, 0)

	rows := pgxmock.NewRows([]string{"category", "price", "time"}).
		AddRow("AppStore", "400", dateStart).
		AddRow("AppStore", "315", dateStart).
		AddRow("AWS", "2700", dateStart.AddDate(0, 0, 1)).
		AddRow("Sport", "980", dateStart.AddDate(0, 0, 1)).
		AddRow("AppStore", "900", dateStart.AddDate(0, 0, 2))

	mock.ExpectQuery(`SELECT category, price, time FROM epxenses`).
		WithArgs(int64(100), dateStart, dateEnd).
		WillReturnRows(rows)

	expenses, err := storage.Get(ctx, entity.UserID(100), dateStart, dateEnd)
	assert.NoError(t, err)

	assert.Equal(t, []entity.Expense{
		entity.NewExpense("AppStore", decimal.New(400, 0), dateStart),
		entity.NewExpense("AppStore", decimal.New(315, 0), dateStart),
		entity.NewExpense("AWS", decimal.New(2700, 0), dateStart.AddDate(0, 0, 1)),
		entity.NewExpense("Sport", decimal.New(980, 0), dateStart.AddDate(0, 0, 1)),
		entity.NewExpense("AppStore", decimal.New(900, 0), dateStart.AddDate(0, 0, 2)),
	}, expenses)
}
