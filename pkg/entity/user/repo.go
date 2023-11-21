package user

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
)

type Repo interface {
	Load(ctx context.Context, userID int64) (*model.User, error)
	Save(ctx context.Context, userID int64, user *model.User) error
}
