package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Expense struct {
	category string
	price    decimal.Decimal
	date     time.Time
}

func NewExpense(category string, price decimal.Decimal, date time.Time) Expense {
	return Expense{
		category: category,
		price:    price,
		date:     date,
	}
}

func (e *Expense) GetCategory() string {
	return e.category
}

func (e *Expense) GetPrice() decimal.Decimal {
	return e.price
}

func (e *Expense) GetDate() time.Time {
	return e.date
}
