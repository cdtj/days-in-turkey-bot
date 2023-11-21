package country

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
	"golang.org/x/text/language"
)

type Service interface {
	CountryInfo(ctx context.Context, language language.Tag, c *model.Country) string
	DefaultCountry(ctx context.Context) *model.Country
	CustomCountry(ctx context.Context, daysCont, daysLimit, resetInterval int) *model.Country
}
