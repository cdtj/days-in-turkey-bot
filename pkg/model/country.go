package model

type Country struct {
	Code          string
	Name          string
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

var (
	defaultCountryCode = "RU"
)

func DefaultCountryCode() string {
	return defaultCountryCode
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

func CustomCountry(cont, limit, reset int) *Country {
	return &Country{
		Code:          "",
		Flag:          "❔",
		DaysContinual: cont,
		DaysLimit:     limit,
		ResetInterval: reset,
	}
}

func (c *Country) GetResetInterval() int {
	return c.ResetInterval
}

func (c *Country) GetDaysCont() int {
	return c.DaysContinual
}

func (c *Country) GetDaysLimit() int {
	return c.DaysLimit
}

func (c *Country) GetCode() string {
	return c.Code
}

func (c *Country) GetFlag() string {
	return c.Flag
}

func (c *Country) GetName() string {
	return c.Name
}
