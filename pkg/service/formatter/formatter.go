package formatter

import (
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/l10n"
)

type Formatter interface {
	TripTree(locale *l10n.Locale, tree *model.TripTree) string
	User(locale *l10n.Locale, user *model.User) string
	Country(locale *l10n.Locale, country *model.Country) string
	FormatMessage(locale *l10n.Locale, messageID string) string
}
