package model

import (
	"cdtj.io/days-in-turkey-bot/service/i18n"
	"golang.org/x/text/language"
)

type User struct {
	Lang        string
	CountryCode string
	Country     Country
	langTag     language.Tag
}

func NewUserConfig(lang, countryCode string) *User {
	return &User{
		Lang:        lang,
		CountryCode: countryCode,
	}
}

func DefaultUser() *User {
	return NewUserConfig(i18n.DefaultLang().String(), DefaultCountryCode())
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
	return u.Lang
}

func (u *User) GetLangTag() language.Tag {
	return u.langTag
}

func (u *User) SetLangTag(tag language.Tag) {
	u.langTag = tag
}
