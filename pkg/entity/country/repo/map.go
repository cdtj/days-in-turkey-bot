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

func (r *CountryRepo) Load(ctx context.Context, userID uint64) (*model.Country, error) {
	u, err := r.db.Load(ctx, userID)
	if err != nil {
		return nil, err
	}
	return u.(*model.Country), err
}

func (r *CountryRepo) Save(ctx context.Context, userID uint64, user *model.Country) error {
	return r.db.Save(ctx, userID, user)
}
