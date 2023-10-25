package user

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
)

type Repo interface {
	Load(ctx context.Context, userID uint64) (*model.User, error)
	Save(ctx context.Context, userID uint64, user *model.User) error
}
