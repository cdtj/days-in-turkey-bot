package service

import (
	"context"

	"cdtj.io/days-in-turkey-bot/entity/country"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"golang.org/x/text/language"
)

var _ country.Service = NewCountryService(nil, nil)

type CountryService struct {
	defaultCountry *model.Country
	fmtr           formatter.Formatter
}

func NewCountryService(fmtr formatter.Formatter, defaultCountry *model.Country) *CountryService {
	return &CountryService{
		fmtr:           fmtr,
		defaultCountry: defaultCountry,
	}
}

func (s *CountryService) CountryInfo(ctx context.Context, language language.Tag, country *model.Country) string {
	return s.fmtr.Country(language, country)
}

func (s *CountryService) DefaultCountry(context.Context) *model.Country {
	return s.defaultCountry
}

func (s *CountryService) CustomCountry(ctx context.Context, daysCont, daysLimit, resetInterval int) *model.Country {
	return &model.Country{
		Code:          "CUSTOM",
		Name:          "",
		Flag:          "📝",
		DaysContinual: daysCont,
		DaysLimit:     daysLimit,
		ResetInterval: resetInterval,
		VisaFree:      true,
	}
}
