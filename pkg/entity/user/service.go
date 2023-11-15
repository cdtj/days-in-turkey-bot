package user

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
	"golang.org/x/text/language"
)

type Service interface {
	UserInfo(ctx context.Context, language language.Tag, user *model.User) string
	CalculateTrip(ctx context.Context, language language.Tag, input string, daysCont, daysLimit, resetInterval int) (string, error)
	DefaultUser(ctx context.Context, userID int64) *model.User
	ParseLanguage(ctx context.Context, languageCode string) language.Tag
}
