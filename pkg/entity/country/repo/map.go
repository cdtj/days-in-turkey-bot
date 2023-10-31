package repo

import (
	"context"

	"cdtj.io/days-in-turkey-bot/db"
	"cdtj.io/days-in-turkey-bot/model"
)

type CountryRepo struct {
	db db.Database
}

func NewCountryRepo(db db.Database) *CountryRepo {
	return &CountryRepo{
		db: db,
	}
}

func (r *CountryRepo) Get(ctx context.Context, userID string) (*model.Country, error) {
	u, err := r.db.Load(ctx, userID)
	if err != nil {
		return nil, err
	}
	return u.(*model.Country), err
}

func (r *CountryRepo) Set(ctx context.Context, userID string, user *model.Country) error {
	return r.db.Save(ctx, userID, user)
}
