package expensepgsqlstorage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

type PgxIface interface {
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
}

type ExpensePgsqlStorage struct {
	conn PgxIface
}

func New(conn PgxIface) *ExpensePgsqlStorage {
	return &ExpensePgsqlStorage{conn: conn}
}

func (s *ExpensePgsqlStorage) Create(ctx context.Context, userID entity.UserID, expense entity.Expense) error {
	_, err := s.conn.Exec(ctx,
		`INSERT INTO epxenses (user_id, category, price, time) VALUES ($1, $2, $3, $4)`,
		int64(userID), expense.GetCategory(), expense.GetPrice().String(), expense.GetDate())

	return errors.Wrap(err, "ExpensePgsqlStorage.Create")
}

func (s *ExpensePgsqlStorage) Get(ctx context.Context, userID entity.UserID, dateStart time.Time, dateEnd time.Time) (
	[]entity.Expense, error,
) {
	rows, err := s.conn.Query(ctx,
		`SELECT category, price, time FROM epxenses
		WHERE user_id = $1 AND time >= $2 AND time < $3
		ORDER BY category`,
		int64(userID), dateStart, dateEnd)
	if err != nil {
		return nil, errors.Wrap(err, "ExpensePgsqlStorage.Get")
	}

	expenses := make([]entity.Expense, 0, rows.CommandTag().RowsAffected())

	var (
		category string
		priceStr string
		date     time.Time
	)

	_, err = pgx.ForEachRow(rows, []any{&category, &priceStr, &date}, func() error {
		price, err := decimal.NewFromString(priceStr)

		if err != nil {
			return errors.Wrap(err, "ExpensePgsqlStorage.Get")
		}

		expenses = append(expenses, entity.NewExpense(category, price, date))

		return nil
	})

	return expenses, errors.Wrap(err, "ExpensePgsqlStorage.Get")
}
