package db

import "context"

type Database interface {
	Keys(ctx context.Context) ([]string, error)
	Load(ctx context.Context, id string) (interface{}, error)
	Save(ctx context.Context, id string, intfc interface{}) error
}
