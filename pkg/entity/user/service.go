package user

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
	"golang.org/x/text/language"
)

type Service interface {
	UserInfo(ctx context.Context, u *model.User) string
	CalculateTrip(ctx context.Context, input string, daysLimit, daysCont, resetInterval int) (string, error)
	LangLookup(ctx context.Context, lang string) (language.Tag, error)
}
