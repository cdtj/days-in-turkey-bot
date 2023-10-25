package db

import (
	"context"
	"errors"
	"sync"
)

type MapDB struct {
	data map[uint64]interface{}
	mu   sync.RWMutex
}

var (
	ErrMapDBNotFound = errors.New("entry not found")
)

func (db *MapDB) Load(ctx context.Context, id uint64) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if data, ok := db.data[id]; ok {
		return data, nil
	}
	return nil, ErrMapDBNotFound
}

func (db *MapDB) Save(ctx context.Context, id uint64, intfc interface{}) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[id] = intfc
	return nil
}
