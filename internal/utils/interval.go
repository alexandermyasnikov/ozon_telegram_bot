package utils

import "time"

const (
	DayInterval   = 1
	WeekInterval  = 2
	MonthInterval = 3

	DaysOfWeek = 7
)

func TruncDate(date time.Time) time.Time {
	y, m, d := date.Date()

	return time.Date(y, m, d, 0, 0, 0, 0, date.Location())
}

func GetInterval(date time.Time, intervalType int) (time.Time, time.Time) {
	date = TruncDate(date)

	switch intervalType {
	case DayInterval:
		return date, date.AddDate(0, 0, 1)
	case WeekInterval:
		offsetToStart := int(date.Weekday() - time.Monday)
		start := date.AddDate(0, 0, -offsetToStart)
		end := start.AddDate(0, 0, DaysOfWeek)

		return start, end
	case MonthInterval:
		offsetToStart := date.Day()
		start := date.AddDate(0, 0, -offsetToStart)
		end := start.AddDate(0, 1, 0)

		return start, end
	default:
		return date, date
	}
}

func IntervalFromStr(interval string) (int, bool) {
	switch interval {
	case "день":
		return DayInterval, true
	case "неделя":
		return WeekInterval, true
	case "месяц":
		return MonthInterval, true
	default:
		return 0, false
	}
}

func IntervalToStr(interval int) (string, bool) {
	switch interval {
	case DayInterval:
		return "день", true
	case WeekInterval:
		return "неделя", true
	case MonthInterval:
		return "месяц", true
	default:
		return "", false
	}
}
