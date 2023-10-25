package db

import (
	"context"
	"errors"
	"sync"
)

type MapDB struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func NewMapDB() *MapDB {
	return &MapDB{
		data: make(map[string]interface{}, 0),
		mu:   sync.RWMutex{},
	}
}

var (
	ErrMapDBNotFound = errors.New("entry not found")
)

func (db *MapDB) Load(ctx context.Context, id string) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if data, ok := db.data[id]; ok {
		return data, nil
	}
	return nil, ErrMapDBNotFound
}

func (db *MapDB) Save(ctx context.Context, id string, intfc interface{}) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[id] = intfc
	return nil
}
