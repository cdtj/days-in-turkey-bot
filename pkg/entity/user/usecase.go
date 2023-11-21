package user

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
	"golang.org/x/text/language"
)

type Usecase interface {
	Create(ctx context.Context, userID int64, lang string) error
	Get(ctx context.Context, userID int64) (*model.User, error)

	GetInfo(ctx context.Context, user *model.User) string
	CalculateTrip(ctx context.Context, user *model.User, datesInput string) (string, error)
	GetLanguage(ctx context.Context, user *model.User) language.Tag

	UpdateLanguage(ctx context.Context, user *model.User, lang string) error
	UpdateCountry(ctx context.Context, user *model.User, country *model.Country) error
}
