package repo

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
)

type CountryDatabase interface {
	Keys(ctx context.Context) ([]string, error)
	Load(ctx context.Context, id string) (interface{}, error)
	Save(ctx context.Context, id string, intfc interface{}) error
}

type CountryRepo struct {
	db CountryDatabase
}

func NewCountryRepo(db CountryDatabase) *CountryRepo {
	return &CountryRepo{
		db: db,
	}
}

func (r *CountryRepo) Get(ctx context.Context, countryID string) (*model.Country, error) {
	u, err := r.db.Load(ctx, countryID)
	if err != nil {
		return nil, err
	}
	return u.(*model.Country), err
}

func (r *CountryRepo) Set(ctx context.Context, countryID string, country *model.Country) error {
	return r.db.Save(ctx, countryID, country)
}

func (r *CountryRepo) Keys(ctx context.Context) ([]string, error) {
	return r.db.Keys(ctx)
}
