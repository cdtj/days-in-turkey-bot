package service

import (
	"context"

	"cdtj.io/days-in-turkey-bot/entity/country"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/l10n"
)

var _ country.Service = NewCountryService(nil)

type CountryService struct {
	fmtr formatter.Formatter
}

func NewCountryService(fmtr formatter.Formatter) *CountryService {
	return &CountryService{
		fmtr: fmtr,
	}
}

func (s *CountryService) Info(ctx context.Context, l *l10n.Locale, c *model.Country) string {
	return s.fmtr.Country(l, c)
}
