package repo

import (
	"context"

	"cdtj.io/days-in-turkey-bot/db"
	"cdtj.io/days-in-turkey-bot/model"
)

type UserRepo struct {
	db db.Database
}

func NewUserRepo(db db.Database) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Load(ctx context.Context, userID string) (*model.User, error) {
	u, err := r.db.Load(ctx, userID)
	if err != nil {
		return nil, err
	}
	return u.(*model.User), err
}
func (r *UserRepo) Save(ctx context.Context, userID string, user *model.User) error {
	return r.db.Save(ctx, userID, user)
}
