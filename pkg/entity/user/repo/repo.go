package repo

import (
	"context"
	"errors"

	"cdtj.io/days-in-turkey-bot/db"
	"cdtj.io/days-in-turkey-bot/entity/user"
	"cdtj.io/days-in-turkey-bot/model"
)

type UserDatabase interface {
	Load(ctx context.Context, id interface{}) (interface{}, error)
	Save(ctx context.Context, id interface{}, intfc interface{}) error
}

type UserRepo struct {
	db UserDatabase
}

var _ = NewUserRepo(nil)

func NewUserRepo(db UserDatabase) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Load(ctx context.Context, userID int64) (*model.User, error) {
	u, err := r.db.Load(ctx, userID)
	if err != nil {
		if errors.Is(err, db.ErrDBEntryNotFound) {
			return nil, user.ErrRepoUserNotFound
		}
		return nil, err
	}
	return u.(*model.User), err
}

func (r *UserRepo) Save(ctx context.Context, userID int64, user *model.User) error {
	return r.db.Save(ctx, userID, user)
}
