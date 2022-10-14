package entity

import "time"

type DateTime struct {
	dateTime time.Time
}

func NewDateTime(year, month, day, hour, min, sec int) DateTime {
	dateTime := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC)

	return DateTime{
		dateTime: dateTime,
	}
}

func NewDateTimeFromTime(t time.Time) DateTime {
	dateTime := t.Truncate(time.Second)

	return DateTime{
		dateTime: dateTime,
	}
}

func (d DateTime) ToTime() time.Time {
	return d.dateTime
}

// ----

type Date struct {
	date time.Time
}

func NewDate(year, moth, day int) Date {
	date := time.Date(year, time.Month(moth), day, 0, 0, 0, 0, time.UTC)

	return Date{
		date: date,
	}
}

func NewDateFromTime(t time.Time) Date {
	y, m, d := t.Date()
	date := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return Date{
		date: date,
	}
}

func (d Date) ToTime() time.Time {
	return d.date
}

func (d Date) ToInt64() int64 {
	return d.date.Unix()
}

func (d *Date) AddDays(days int) {
	d.date = d.date.AddDate(0, 0, days)
}
