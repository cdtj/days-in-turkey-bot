package user

import "golang.org/x/text/language"

type UserConfig struct {
	lang    language.Tag
	country *UserCountry
}

func NewUserConfig(lang language.Tag, country *UserCountry) *UserConfig {
	return &UserConfig{
		lang:    lang,
		country: country,
	}
}

func (u *UserConfig) GetDaysReset() int {
	return u.country.daysReset
}

type UserCountry struct {
	code          string
	daysContinual int
	daysLimit     int
	daysReset     int
}

func NewUserCountry(code string, daysContinual, daysLimit, daysMustPassed int) *UserCountry {
	return &UserCountry{
		code:          code,
		daysContinual: daysContinual,
		daysLimit:     daysLimit,
		daysReset:     daysMustPassed,
	}
}
