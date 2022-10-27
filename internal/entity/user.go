package entity

import "github.com/shopspring/decimal"

type UserID int64

type User struct {
	id              UserID
	defaultCurrency string
	dayLimit        decimal.Decimal
	weekLimit       decimal.Decimal
	monthLimit      decimal.Decimal
}

func NewUser(userID UserID) User {
	return User{
		id:              userID,
		defaultCurrency: "",
		dayLimit:        decimal.Zero,
		weekLimit:       decimal.Zero,
		monthLimit:      decimal.Zero,
	}
}

func (u User) GetID() UserID {
	return u.id
}

func (u User) GetDefaultCurrency() string {
	return u.defaultCurrency
}

func (u *User) SetDefaultCurrency(currency string) {
	u.defaultCurrency = currency
}

func (u User) GetDayLimit() decimal.Decimal {
	return u.dayLimit
}

func (u *User) SetDayLimit(limit decimal.Decimal) {
	u.dayLimit = limit
}

func (u User) GetWeekLimit() decimal.Decimal {
	return u.weekLimit
}

func (u *User) SetWeekLimit(limit decimal.Decimal) {
	u.weekLimit = limit
}

func (u User) GetMonthLimit() decimal.Decimal {
	return u.monthLimit
}

func (u *User) SetMonthLimit(limit decimal.Decimal) {
	u.monthLimit = limit
}
