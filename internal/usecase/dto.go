package usecase

import (
	"time"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

type GetAllCurrencyNamesRespDTO struct {
	Currencies []string
}

type SetDefaultCurrencyReqDTO struct {
	UserID   int64
	Currency string
}

type AddExpenseReqDTO struct {
	UserID   int64
	Category string
	Price    float64
	Date     time.Time
	Currency string
}

type GetReportReqDTO struct {
	UserID   int64
	Date     time.Time
	Days     int
	Currency string
}

type GetReportRespDTO struct {
	Currency   string
	Categories map[string]entity.Decimal // TODO заменить на slice, порядок может меняться в тестах
}

type UpdateCurrencyReqDTO struct {
	Currency string
	Rate     float64
}
