package usecase

import (
	"context"

	"cdtj.io/days-in-turkey-bot/entity/country"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/l10n"
)

var _ country.Usecase = NewCountryUsecase(nil, nil)

type CountryUsecase struct {
	repo    country.Repo
	service country.Service
}

func NewCountryUsecase(repo country.Repo, service country.Service) *CountryUsecase {
	return &CountryUsecase{
		repo:    repo,
		service: service,
	}
}

func (u *CountryUsecase) Get(ctx context.Context, countryID string) (*model.Country, error) {
	return u.repo.Get(ctx, countryID)
}

func (u *CountryUsecase) Set(ctx context.Context, countryID string, country *model.Country) error {
	return u.repo.Set(ctx, countryID, country)
}

func (u *CountryUsecase) Info(ctx context.Context, countryID string) (string, error) {
	c, err := u.Get(ctx, countryID)
	if err != nil {
		return "", err
	}
	// default locale is enough for debugging
	return u.service.Info(ctx, l10n.NewLocale(l10n.DefaultLang()), c), nil
}
