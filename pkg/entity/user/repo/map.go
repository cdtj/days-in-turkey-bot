package repo

import (
	"context"
	"errors"

	"cdtj.io/days-in-turkey-bot/db"
	"cdtj.io/days-in-turkey-bot/entity/user"
	"cdtj.io/days-in-turkey-bot/model"
)

type UseryDatabase interface {
	Load(ctx context.Context, id string) (interface{}, error)
	Save(ctx context.Context, id string, intfc interface{}) error
}

type UserRepo struct {
	db UseryDatabase
}

func NewUserRepo(db UseryDatabase) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Load(ctx context.Context, userID string) (*model.User, error) {
	u, err := r.db.Load(ctx, userID)
	if err != nil {
		if errors.Is(err, db.ErrDBEntryNotFound) {
			return nil, user.ErrRepoUserNotFound
		}
		return nil, err
	}
	return u.(*model.User), err
}

func (r *UserRepo) Save(ctx context.Context, userID string, user *model.User) error {
	return r.db.Save(ctx, userID, user)
}
