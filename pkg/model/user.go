package model

import (
	"golang.org/x/text/language"
)

type User struct {
	ID       int64
	Country  Country
	Language language.Tag
}

func NewUser(id int64, language language.Tag, country Country) *User {
	return &User{
		ID:       id,
		Country:  country,
		Language: language,
	}
}

func (u *User) GetID() int64 {
	return u.ID
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

func (u *User) GetLanguage() language.Tag {
	return u.Language
}

func (u *User) SetLanguage(language language.Tag) {
	u.Language = language
}

func (u *User) SetCountry(country Country) {
	u.Country = country
}
