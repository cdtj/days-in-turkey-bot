package country

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
)

type Repo interface {
	Load(ctx context.Context, countryID string) (*model.Country, error)
	Save(ctx context.Context, countryID string, country *model.Country) error
}
