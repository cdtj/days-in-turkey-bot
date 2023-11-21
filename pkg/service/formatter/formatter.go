package formatter

import (
	"cdtj.io/days-in-turkey-bot/model"
	"golang.org/x/text/language"
)

type Formatter interface {
	TripTree(language language.Tag, tree *model.TripTree) string
	User(language language.Tag, user *model.User) string
	Country(language language.Tag, country *model.Country) string
	FormatMessage(language language.Tag, messageID string) string

	Welcome(language language.Tag) string
	TripExplanation(language language.Tag) string
}
