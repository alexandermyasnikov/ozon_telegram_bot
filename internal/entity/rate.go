package entity

type Rate struct {
	code  string
	ratio Decimal
	time  DateTime
}

func NewRate(code string, ratio Decimal, time DateTime) Rate {
	return Rate{
		code:  code,
		ratio: ratio,
		time:  time,
	}
}

func (r Rate) GetCode() string {
	return r.code
}

func (r Rate) GetRatio() Decimal {
	return r.ratio
}

func (r Rate) GetTime() DateTime {
	return r.time
}
