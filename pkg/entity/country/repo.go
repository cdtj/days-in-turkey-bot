package country

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
)

type Repo interface {
	Keys(ctx context.Context) ([]string, error)
	Cache(ctx context.Context) []*model.Country
	BuildCache(ctx context.Context) error
	Get(ctx context.Context, countryID string) (*model.Country, error)
	Set(ctx context.Context, countryID string, country *model.Country) error
}
