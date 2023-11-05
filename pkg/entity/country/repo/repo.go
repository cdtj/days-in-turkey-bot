package repo

import (
	"context"
	"sort"
	"sync"

	"cdtj.io/days-in-turkey-bot/model"
)

type CountryDatabase interface {
	Keys(ctx context.Context) ([]interface{}, error)
	Load(ctx context.Context, id interface{}) (interface{}, error)
	Save(ctx context.Context, id interface{}, intfc interface{}) error
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
	if err := constructor(context.Background(), repo); err != nil {
		return nil
	}
	return repo
}

func constructor(ctx context.Context, repo *CountryRepo) error {
	mu := new(sync.Mutex)
	mu.Lock()
	defer mu.Unlock()
	keys, err := repo.Keys(ctx)
	if err != nil {
		return nil
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	for _, key := range keys {
		country, err := repo.Get(ctx, key)
		if err != nil {
			return nil
		}
		repo.cache = append(repo.cache, country)
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
