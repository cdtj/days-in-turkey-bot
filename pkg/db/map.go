package db

import (
	"context"
	"log/slog"
	"sync"
)

type MapDB struct {
	data map[interface{}]interface{}
	mu   sync.RWMutex
}

func NewMapDB() *MapDB {
	return &MapDB{
		data: make(map[interface{}]interface{}, 0),
		mu:   sync.RWMutex{},
	}
}

func (db *MapDB) Load(ctx context.Context, id any) (any, error) {
	slog.Debug("loading", "id", id)
	db.mu.RLock()
	defer db.mu.RUnlock()
	if data, ok := db.data[id]; ok {
		return data, nil
	}
	return nil, ErrDBEntryNotFound
}

func (db *MapDB) Save(ctx context.Context, id any, data any) error {
	slog.Debug("saving", "id", id)
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[id] = data
	return nil
}

func (db *MapDB) Keys(ctx context.Context) ([]any, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	keys := make([]any, 0, len(db.data))
	for k := range db.data {
		keys = append(keys, k)
	}

	return keys, nil
}
