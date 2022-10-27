package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Rate struct {
	code  string
	ratio decimal.Decimal
	time  time.Time
}

func NewRate(code string, ratio decimal.Decimal, time time.Time) Rate {
	return Rate{
		code:  code,
		ratio: ratio,
		time:  time,
	}
}

func (r Rate) GetCode() string {
	return r.code
}

func (r Rate) GetRatio() decimal.Decimal {
	return r.ratio
}

func (r Rate) GetTime() time.Time {
	return r.time
}
