package db

import "context"

type Database interface {
	Load(ctx context.Context, id string) (interface{}, error)
	Save(ctx context.Context, id string, intfc interface{}) error
}
