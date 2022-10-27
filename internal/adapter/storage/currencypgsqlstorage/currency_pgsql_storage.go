package currencypgsqlstorage

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

type CurrencyPgsqlStorage struct {
	conn PgxIface
}

func New(conn PgxIface) *CurrencyPgsqlStorage {
	return &CurrencyPgsqlStorage{conn: conn}
}

func (s *CurrencyPgsqlStorage) Get(ctx context.Context, currency string) (entity.Rate, error) {
	var (
		code     string
		ratioStr string
		date     time.Time
	)

	err := s.conn.QueryRow(ctx,
		`SELECT code, ratio, time FROM currencies WHERE code = $1`,
		currency).Scan(&code, &ratioStr, &date)
	if err != nil {
		return entity.Rate{}, errors.Wrap(err, "CurrencyPgsqlStorage.Get")
	}

	ratio, err := decimal.NewFromString(ratioStr)
	if err != nil {
		return entity.Rate{}, errors.Wrap(err, "CurrencyPgsqlStorage.Get")
	}

	return entity.NewRate(code, ratio, date), errors.Wrap(err, "CurrencyPgsqlStorage.Get")
}

func (s *CurrencyPgsqlStorage) GetAll(ctx context.Context) ([]entity.Rate, error) {
	rows, err := s.conn.Query(ctx,
		`SELECT code, ratio, time FROM currencies`)
	if err != nil {
		return nil, errors.Wrap(err, "CurrencyPgsqlStorage.GetAll")
	}

	var (
		code     string
		ratioStr string
		date     time.Time

		rates []entity.Rate
	)

	_, err = pgx.ForEachRow(rows, []any{&code, &ratioStr, &date}, func() error {
		ratio, err := decimal.NewFromString(ratioStr)
		if err != nil {
			return errors.Wrap(err, "CurrencyPgsqlStorage.GetAll")
		}

		rates = append(rates, entity.NewRate(code, ratio, date))

		return nil
	})

	return rates, errors.Wrap(err, "CurrencyPgsqlStorage.GetAll")
}

func (s *CurrencyPgsqlStorage) Update(ctx context.Context, rate entity.Rate) error {
	_, err := s.conn.Exec(ctx,
		`INSERT INTO currencies (code, ratio, time) VALUES ($1, $2, $3)
		ON CONFLICT (code) DO UPDATE
		SET ratio = EXCLUDED.ratio, time = EXCLUDED.time`,
		rate.GetCode(), rate.GetRatio().String(), rate.GetTime())

	return errors.Wrap(err, "CurrencyPgsqlStorage.Update")
}
