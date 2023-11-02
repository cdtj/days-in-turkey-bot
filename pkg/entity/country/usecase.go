package country

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
)

type Usecase interface {
	Keys(ctx context.Context) ([]string, error)
	Get(ctx context.Context, countryID string) (*model.Country, error)
	Set(ctx context.Context, countryID string, country *model.Country) error
	Info(ctx context.Context, userID string) (string, error)
}
