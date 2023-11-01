package model

import (
	"cdtj.io/days-in-turkey-bot/service/l10n"
	"golang.org/x/text/language"
)

type User struct {
	Country *Country
	locale  *l10n.Locale
}

func NewUserConfig(lang language.Tag, country *Country) *User {
	return &User{
		locale:  l10n.NewLocale(lang),
		Country: country,
	}
}

func DefaultUser() *User {
	return NewUserConfig(l10n.DefaultLang(), DefaultCountry())
}

func (u *User) SetLocale(lang language.Tag) {
	u.locale = l10n.NewLocale(lang)
}

func (u *User) GetLocale() *l10n.Locale {
	return u.locale
}

func (u *User) GetResetInterval() int {
	return u.Country.GetResetInterval()
}

func (u *User) GetDaysCont() int {
	return u.Country.GetDaysCont()
}

func (u *User) GetDaysLimit() int {
	return u.Country.GetDaysLimit()
}

func (u *User) GetLang() string {
	return u.locale.Lang()
}
