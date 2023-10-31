package model

import "fmt"

type Country struct {
	Code          string
	Flag          string
	DaysContinual int
	DaysLimit     int
	ResetInterval int
}

func NewCountry(code, flag string, daysContinual, daysLimit, resetInterval int) *Country {
	return &Country{
		Code:          code,
		Flag:          flag,
		DaysContinual: daysContinual,
		DaysLimit:     daysLimit,
		ResetInterval: resetInterval,
	}
}

func DefaultCountry() *Country {
	return &Country{
		Code:          "RU",
		Flag:          "🇷🇺",
		DaysContinual: 60,
		DaysLimit:     90,
		ResetInterval: 180,
	}
}

func (c *Country) String() string {
	return fmt.Sprintf("[%s %s]\nDays(Continual: %d, Limit: %d, Reset: %d)", c.Code, c.Flag, c.DaysContinual, c.DaysLimit, c.ResetInterval)
}
