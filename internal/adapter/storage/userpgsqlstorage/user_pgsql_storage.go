package userpgsqlstorage

import (
	"context"

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

type UserPgsqlStorage struct {
	conn PgxIface
}

func New(conn PgxIface) *UserPgsqlStorage {
	return &UserPgsqlStorage{conn: conn}
}

func (s *UserPgsqlStorage) GetDefaultCurrency(ctx context.Context, userID entity.UserID) (string, error) {
	var currency string

	err := s.conn.QueryRow(ctx,
		`SELECT currency FROM users WHERE id = $1`,
		int64(userID)).Scan(&currency)
	if err != nil {
		return "", errors.Wrap(err, "UserPgsqlStorage.GetDefaultCurrency")
	}

	return currency, errors.Wrap(err, "UserPgsqlStorage.GetDefaultCurrency")
}

func (s *UserPgsqlStorage) UpdateDefaultCurrency(ctx context.Context, userID entity.UserID, currency string) error {
	_, err := s.conn.Exec(ctx,
		`INSERT INTO users (id, currency) VALUES ($1, $2)
			ON CONFLICT (id) DO UPDATE SET currency = $2`,
		int64(userID), currency)

	return errors.Wrap(err, "UpdateDefaultCurrency.UpdateDefaultCurrency")
}

func (s *UserPgsqlStorage) GetLimits(ctx context.Context, userID entity.UserID) (
	decimal.Decimal, decimal.Decimal, decimal.Decimal, error,
) {
	var (
		dayLimitStr   string
		weekLimitStr  string
		monthLimitStr string
	)

	err := s.conn.QueryRow(ctx,
		`SELECT day_limit, week_limit, month_limit FROM users WHERE id = $1`,
		int64(userID)).Scan(&dayLimitStr, &weekLimitStr, &monthLimitStr)
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, decimal.Decimal{}, errors.Wrap(err, "UserPgsqlStorage.GetLimits")
	}

	dayLimit, err := decimal.NewFromString(dayLimitStr)
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, decimal.Decimal{}, errors.Wrap(err, "UserPgsqlStorage.GetLimits")
	}

	weekLimit, err := decimal.NewFromString(weekLimitStr)
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, decimal.Decimal{}, errors.Wrap(err, "UserPgsqlStorage.GetLimits")
	}

	monthLimit, err := decimal.NewFromString(monthLimitStr)
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, decimal.Decimal{}, errors.Wrap(err, "UserPgsqlStorage.GetLimits")
	}

	return dayLimit, weekLimit, monthLimit, errors.Wrap(err, "UserPgsqlStorage.GetLimits")
}

func (s *UserPgsqlStorage) UpdateDayLimit(ctx context.Context, userID entity.UserID, limit decimal.Decimal) error {
	_, err := s.conn.Exec(ctx,
		`INSERT INTO users (id, day_limit) VALUES ($1, $2)
			ON CONFLICT (id) DO UPDATE SET day_limit = $2`,
		int64(userID), limit.String())

	return errors.Wrap(err, "UserPgsqlStorage.UpdateDayLimit")
}

func (s *UserPgsqlStorage) UpdateWeekLimit(ctx context.Context, userID entity.UserID, limit decimal.Decimal) error {
	_, err := s.conn.Exec(ctx,
		`INSERT INTO users (id, week_limit) VALUES ($1, $2)
			ON CONFLICT (id) DO UPDATE SET week_limit = $2`,
		int64(userID), limit.String())

	return errors.Wrap(err, "UserPgsqlStorage.UpdateWeekLimit")
}

func (s *UserPgsqlStorage) UpdateMonthLimit(ctx context.Context, userID entity.UserID, limit decimal.Decimal) error {
	_, err := s.conn.Exec(ctx,
		`INSERT INTO users (id, month_limit) VALUES ($1, $2)
			ON CONFLICT (id) DO UPDATE SET month_limit = $2`,
		int64(userID), limit.String())

	return errors.Wrap(err, "UserPgsqlStorage.UpdateMonthLimit")
}
