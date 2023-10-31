package user

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
)

type Usecase interface {
	Get(ctx context.Context, userID string) (*model.User, error)
	Info(ctx context.Context, userID string) (string, error)
	CalculateTrip(ctx context.Context, userID string, input string) (string, error)
	UpdateLang(ctx context.Context, userID string, lang string) (string, error)
	UpdateCountry(ctx context.Context, userID string, countryID string) (string, error)
}
