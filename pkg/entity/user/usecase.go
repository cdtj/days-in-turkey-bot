package user

import (
	"context"

	"golang.org/x/text/language"
)

type Usecase interface {
	Create(ctx context.Context, userID, lang string) error
	Info(ctx context.Context, userID string) (string, error)
	CalculateTrip(ctx context.Context, userID string, input string) (string, error)
	UpdateLang(ctx context.Context, userID string, lang string) error
	UpdateCountry(ctx context.Context, userID string, countryID string) error
	GetLang(ctx context.Context, userID string) language.Tag
}
