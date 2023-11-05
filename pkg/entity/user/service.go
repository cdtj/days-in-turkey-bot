package user

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/i18n"
	"golang.org/x/text/language"
)

type Service interface {
	UserInfo(ctx context.Context, l *i18n.Locale, u *model.User) string
	CalculateTrip(ctx context.Context, l *i18n.Locale, input string, daysLimit, daysCont, resetInterval int) (string, error)
	LangLookup(ctx context.Context, lang string) (language.Tag, error)
}
