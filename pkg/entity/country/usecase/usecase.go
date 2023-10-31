package usecase

import (
	"context"

	"cdtj.io/days-in-turkey-bot/entity/country"
	"cdtj.io/days-in-turkey-bot/model"
)

type CountryUsecase struct {
	repo country.Repo
}

func NewCountryUsecase(repo country.Repo) *CountryUsecase {
	return &CountryUsecase{
		repo: repo,
	}
}

func (u *CountryUsecase) Get(ctx context.Context, countryID string) (*model.Country, error) {
	return u.repo.Get(ctx, countryID)
}

func (u *CountryUsecase) Set(ctx context.Context, countryID string, country *model.Country) error {
	return u.repo.Set(ctx, countryID, country)
}
