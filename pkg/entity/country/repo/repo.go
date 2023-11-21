package repo

import (
	"context"
	"log/slog"
	"sort"
	"sync"

	"cdtj.io/days-in-turkey-bot/model"
)

type CountryDatabase interface {
	Keys(ctx context.Context) ([]any, error)
	Load(ctx context.Context, id any) (any, error)
	Save(ctx context.Context, id any, intfc any) error
}

type CountryRepo struct {
	db    CountryDatabase
	cache []*model.Country
}

func NewCountryRepo(db CountryDatabase) *CountryRepo {
	repo := &CountryRepo{
		db:    db,
		cache: make([]*model.Country, 0),
	}
	return repo
}

func (r *CountryRepo) BuildCache(ctx context.Context) error {
	mu := new(sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	keys, err := r.Keys(ctx)
	if err != nil {
		return nil
	}
	slog.Debug("constructing countries", "keys", keys)
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	for _, key := range keys {
		country, err := r.Get(ctx, key)
		if err != nil {
			return nil
		}
		slog.Debug("caching country", "key", key, "country", country)
		r.cache = append(r.cache, country)
	}
	return nil
}

func (r *CountryRepo) Get(ctx context.Context, countryID string) (*model.Country, error) {
	u, err := r.db.Load(ctx, countryID)
	if err != nil {
		return nil, err
	}
	return u.(*model.Country), err
}

func (r *CountryRepo) Set(ctx context.Context, countryID string, country *model.Country) error {
	return r.db.Save(ctx, countryID, country)
}

func (r *CountryRepo) Keys(ctx context.Context) ([]string, error) {
	keys, err := r.db.Keys(ctx)
	if err != nil {
		return nil, err
	}
	strKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		strKeys = append(strKeys, key.(string))
	}
	return strKeys, nil
}

func (r *CountryRepo) Cache(ctx context.Context) []*model.Country {
	return r.cache
}
