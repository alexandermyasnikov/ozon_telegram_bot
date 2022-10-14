package entity

import (
	"math"
)

const decimalScaleDefault = 5

type Decimal struct {
	number int64
}

func NewDecimal(number int64, scale uint64) Decimal {
	for scale < decimalScaleDefault {
		number *= 10
		scale++
	}

	for scale > decimalScaleDefault {
		number /= 10
		scale--
	}

	return Decimal{
		number: number,
	}
}

func NewDecimalFromFloat(number float64) Decimal {
	decimal := Decimal{
		number: int64(number * math.Pow10(decimalScaleDefault)),
	}

	return decimal
}

func (d Decimal) ToFloat() float64 {
	return float64(d.number) / math.Pow10(decimalScaleDefault)
}

func (d *Decimal) Add(decimal Decimal) {
	d.number += decimal.number
}

func (d *Decimal) Mult(decimal Decimal) {
	d.number *= decimal.number
	for i := decimalScaleDefault; i > 0; i-- {
		d.number /= 10
	}
}

func (d *Decimal) Div(decimal Decimal) {
	for i := decimalScaleDefault; i > 0; i-- {
		d.number *= 10
	}

	d.number /= decimal.number
}
