package formatter

import (
	"fmt"

	"cdtj.io/days-in-turkey-bot/model"
)

type TelegramFormatter struct {
}

func (f *TelegramFormatter) TripTree(tree *model.TripTree) string {
	result := ""
	for i := tree; i != nil; i = i.Prev {
		result += fmt.Sprintf("trip: %q - %q @ %d / %d\n", i.StartDate, i.EndDate, i.TripDays, i.PeriodDays)
	}
	return result
}

func (f *TelegramFormatter) User(user *model.User) string {
	return ""
}
