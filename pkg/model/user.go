package model

import (
	"golang.org/x/text/language"
)

// User is struct to store User settings,
// fields are exported in case to store them in DB
type User struct {
	ID           int64
	Country      Country
	LanguageCode string
	language     language.Tag `gob:"-"`
}

func NewUser(id int64, languageCode string, country Country) *User {
	return &User{
		ID:           id,
		Country:      country,
		LanguageCode: languageCode,
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
	return u.language
}

func (u *User) SetLanguage(language language.Tag) {
	u.language = language
}

func (u *User) GetLanguageCode() string {
	return u.LanguageCode
}

func (u *User) SetLanguageCode(languageCode string) {
	u.LanguageCode = languageCode
}

func (u *User) SetCountry(country Country) {
	u.Country = country
}
