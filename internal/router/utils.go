package router

import (
	"fmt"
	"strings"
)

const (
	DayToDays   = 1
	WeekToDays  = 7
	MonthToDays = 31
	YearToDays  = 365
)

func genRE(args ...[]string) string {
	argsRE := make([]string, 0, len(args))
	for _, strs := range args {
		argsRE = append(argsRE, fmt.Sprintf("(?P<%s>%s)", strs[0], strings.Join(strs[1:], "|")))
	}

	ans := `^` + strings.Join(argsRE, `(?:\S*)(?:\s+)`)

	return ans
}

func dateToDays(date string) int {
	switch date {
	case "ден":
		return DayToDays
	case "нед":
		return WeekToDays
	case "мес":
		return MonthToDays
	case "год":
		return YearToDays
	default:
		return MonthToDays
	}
}
