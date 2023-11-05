package country

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
	"golang.org/x/text/language"
)

type Usecase interface {
	Get(ctx context.Context, countryID string) (*model.Country, error)
	// Set is commented out because we don't want to affect the list through API,
	// to modify Country/Countries make relevant changes to assets/country,
	// and re-init the CountryMapRepo
	// Set(ctx context.Context, countryID string, country *model.Country) error
	Lookup(ctx context.Context, countryID string, daysCont, daysLimit, resetInterval int) (*model.Country, error)

	// ListFromRepo deprecated
	ListFromRepo(ctx context.Context) ([]*model.Country, error)
	ListFromCache(ctx context.Context) []*model.Country
	GetInfo(ctx context.Context, language language.Tag, country *model.Country) (string, error)

	DefaultCountry(ctx context.Context) *model.Country
}
