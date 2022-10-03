package storage

import (
	"time"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

type ProductStorageInterface interface {
	Add(userID int64, product entity.Product) error
	GetAll(userID int64, date time.Time, days int) ([]entity.Product, error)
}
