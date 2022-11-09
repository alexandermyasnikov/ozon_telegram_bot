package usecase

import (
	"time"

	"github.com/shopspring/decimal"
)

type MessageInfo struct {
	UserID int64     `json:"user_id,omitempty"`
	Date   time.Time `json:"date,omitempty"`
}

type Command struct {
	MessageInfo
	Name                      string                     `json:"name"`
	SetDefaultCurrencyReqDTO  *SetDefaultCurrencyReqDTO  `json:"set_default_currency_req_dto,omitempty"`
	SetDefaultCurrencyRespDTO *SetDefaultCurrencyRespDTO `json:"set_default_currency_resp_dto,omitempty"`
	AddExpenseReqDTO          *AddExpenseReqDTO          `json:"add_expense_req_dto,omitempty"`
	AddExpenseRespDTO         *AddExpenseRespDTO         `json:"add_expense_resp_dto,omitempty"`
	GetReportReqDTO           *GetReportReqDTO           `json:"get_report_req_dto,omitempty"`
	GetReportRespDTO          *GetReportRespDTO          `json:"get_report_resp_dto,omitempty"`
	SetLimitReqDTO            *SetLimitReqDTO            `json:"set_limit_req_dto,omitempty"`
	SetLimitRespDTO           *SetLimitRespDTO           `json:"set_limit_resp_dto,omitempty"`
	GetLimitsReqDTO           *GetLimitsReqDTO           `json:"get_limits_req_dto,omitempty"`
	GetLimitsRespDTO          *GetLimitsRespDTO          `json:"get_limits_resp_dto,omitempty"`
}

type CommandAddExpense struct {
	Req  AddExpenseReqDTO
	Resp AddExpenseRespDTO
}

type SetDefaultCurrencyReqDTO struct {
	UserID   int64
	Currency string
}

type SetDefaultCurrencyRespDTO struct {
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
