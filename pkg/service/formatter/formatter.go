package formatter

import (
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/i18n"
)

type Formatter interface {
	TripTree(locale *i18n.Locale, tree *model.TripTree) string
	User(locale *i18n.Locale, user *model.User) string
	Country(locale *i18n.Locale, country *model.Country) string
}
