package service

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/formatter"
)

type CountryService struct {
	fmtr formatter.Formatter
}

func NewCountryService(fmtr formatter.Formatter) *CountryService {
	return &CountryService{
		fmtr: fmtr,
	}
}

func (s *CountryService) CountryInfo(ctx context.Context, c *model.Country) string {
	return s.fmtr.Country(c)
}
