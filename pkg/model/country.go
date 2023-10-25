package model

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
