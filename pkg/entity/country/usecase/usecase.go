package usecase

import (
	"context"

	"cdtj.io/days-in-turkey-bot/entity/country"
)

type CountryUsecase struct {
	repo country.Repo
}

func NewCountryUsecase(repo country.Repo) *CountryUsecase {
	return &CountryUsecase{
		repo: repo,
	}
}

func (u *CountryUsecase) Create(ctx context.Context) error {
	return nil
}

func (u *CountryUsecase) Update(ctx context.Context) error {
	return nil
}
