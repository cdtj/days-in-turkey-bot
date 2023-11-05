package db

import (
	"context"
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

func (db *MapDB) Load(ctx context.Context, id interface{}) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if data, ok := db.data[id]; ok {
		return data, nil
	}
	return nil, ErrDBEntryNotFound
}

func (db *MapDB) Save(ctx context.Context, id interface{}, intfc interface{}) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[id] = intfc
	return nil
}

func (db *MapDB) Keys(ctx context.Context) ([]interface{}, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	keys := make([]interface{}, 0, len(db.data))
	for k := range db.data {
		keys = append(keys, k)
	}

	return keys, nil
}
