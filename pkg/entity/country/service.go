package country

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
)

type Service interface {
	CountryInfo(ctx context.Context, c *model.Country) string
}
