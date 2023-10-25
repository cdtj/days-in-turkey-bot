package country

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
)

type Repo interface {
	Load(ctx context.Context, countryID uint64) (*model.Country, error)
	Save(ctx context.Context, countryID uint64, user *model.Country) error
}
