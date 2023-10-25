package user

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
)

type Repo interface {
	Load(ctx context.Context, userID string) (*model.User, error)
	Save(ctx context.Context, userID string, user *model.User) error
}
