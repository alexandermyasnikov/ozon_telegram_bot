package usecase

import (
	"time"

	"github.com/shopspring/decimal"
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
	Price    decimal.Decimal
	Date     time.Time
}

type AddExpenseRespDTO struct {
	Limits   map[int]decimal.Decimal
	Currency string
}

type GetReportReqDTO struct {
	UserID       int64
	Date         time.Time
	IntervalType int
}

type GetReportRespDTO struct {
	Currency string
	Expenses []ExpenseReportDTO
}

type UpdateCurrencyReqDTO struct {
	Currency string
	Rate     decimal.Decimal
}

type SetLimitReqDTO struct {
	UserID       int64
	Limit        decimal.Decimal
	IntervalType int
}

type SetLimitRespDTO struct {
	Currency string
}

type GetLimitsReqDTO struct {
	UserID int64
}

type GetLimitsRespDTO struct {
	Limits   map[int]decimal.Decimal
	Currency string
}

// ----

type ExpenseReportDTO struct {
	Category string
	Sum      decimal.Decimal
}
