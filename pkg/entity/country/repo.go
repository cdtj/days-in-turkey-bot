package country

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
)

type Repo interface {
	Get(ctx context.Context, countryID string) (*model.Country, error)
	Set(ctx context.Context, countryID string, country *model.Country) error
}
