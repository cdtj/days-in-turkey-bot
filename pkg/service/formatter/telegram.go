package formatter

import (
	"fmt"

	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/l10n"
	"golang.org/x/text/language"
)

type TelegramFormatter struct {
	lang language.Tag
}

func NewTelegramFormatter() *TelegramFormatter {
	return &TelegramFormatter{}
}

func (f *TelegramFormatter) TripTree(tree *model.TripTree) string {
	result := ""
	for i := tree; i != nil; i = i.Prev {
		result += fmt.Sprintf("%s - %s @ Days: %d of %d\n", l10n.FormatDate(i.StartDate), l10n.FormatDate(i.EndDate), i.TripDays, i.PeriodDays)
	}
	return result
}

func (f *TelegramFormatter) User(user *model.User) string {
	return user.String()
}

func (f *TelegramFormatter) Country(country *model.Country) string {
	return country.String()
}
