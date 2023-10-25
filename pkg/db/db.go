package db

import "context"

type Database interface {
	Load(ctx context.Context, id uint64) (interface{}, error)
	Save(ctx context.Context, id uint64, intfc interface{}) error
}
