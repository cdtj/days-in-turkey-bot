package repo

import (
	"context"
	"errors"
	"log/slog"

	"cdtj.io/days-in-turkey-bot/db"
	"cdtj.io/days-in-turkey-bot/entity/user"
	"cdtj.io/days-in-turkey-bot/model"
)

var _ UserDatabase = NewUserBoltDBAdaptor(db.NewBoltDB("", ""))
var _ UserDatabase = db.NewMapDB()

type UserDatabase interface {
	Load(ctx context.Context, key any) (any, error)
	Save(ctx context.Context, key any, data any) error
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
		slog.Error("user repo", "userID", userID, "err", err)
		if errors.Is(err, db.ErrDBEntryNotFound) || errors.Is(err, db.ErrDBBucketNotFound) {
			return nil, user.ErrRepoUserNotFound
		}
		return nil, err
	}
	return u.(*model.User), err
}

func (r *UserRepo) Save(ctx context.Context, userID int64, user *model.User) error {
	return r.db.Save(ctx, userID, user)
}
