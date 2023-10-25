package user

import "context"

type Usecase interface {
	Create(ctx context.Context, userID uint64) error
	Get(ctx context.Context, userID uint64) error
	Calc(ctx context.Context, userID uint64) error
}
