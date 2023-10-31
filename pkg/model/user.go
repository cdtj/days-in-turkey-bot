package model

import (
	"fmt"

	"golang.org/x/text/language"
)

type User struct {
	Lang    language.Tag
	Country *Country
}

func NewUserConfig(lang language.Tag, country *Country) *User {
	return &User{
		Lang:    lang,
		Country: country,
	}
}

func (u *User) GetResetInterval() int {
	return u.Country.ResetInterval
}

func (u *User) GetDaysCont() int {
	return u.Country.DaysContinual
}

func (u *User) GetDaysLimit() int {
	return u.Country.DaysLimit
}

func DefaultUser() *User {
	return NewUserConfig(language.Russian, DefaultCountry())
}

func (u *User) String() string {
	return fmt.Sprintf("Language: %s\nCountry: %s", u.Lang, u.Country.String())
}
