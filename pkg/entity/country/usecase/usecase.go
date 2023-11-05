package usecase

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"sort"

	"cdtj.io/days-in-turkey-bot/assets"
	"cdtj.io/days-in-turkey-bot/entity/country"
	"cdtj.io/days-in-turkey-bot/model"
	"github.com/BurntSushi/toml"
	"golang.org/x/text/language"
)

var _ country.Usecase = NewCountryUsecase(nil, nil)

const countryFiles = "country"

type CountryUsecase struct {
	repo    country.Repo
	service country.Service
}

func NewCountryUsecase(repo country.Repo, service country.Service) *CountryUsecase {
	uc := &CountryUsecase{
		repo:    repo,
		service: service,
	}
	if err := uc.constructor(); err != nil {
		slog.Error("failed to init CountryUsecase", "err", err)
		return nil
	}
	return uc
}

func (uc *CountryUsecase) constructor() error {
	if uc.repo != nil {
		if err := uc.loadFromFiles(context.Background(), countryFiles); err != nil {
			return err
		}
	}
	return nil
}

func (uc *CountryUsecase) Get(ctx context.Context, countryID string) (*model.Country, error) {
	return uc.repo.Get(ctx, countryID)
}

func (uc *CountryUsecase) Lookup(ctx context.Context, countryID string, daysCont, daysLimit, resetInterval int) (*model.Country, error) {
	if countryID != "" {
		country, err := uc.repo.Get(ctx, countryID)
		if err != nil {
			return nil, err
		}
		return country, nil
	}
	return uc.service.CustomCountry(ctx, daysCont, daysLimit, resetInterval), nil
}

func (uc *CountryUsecase) Set(ctx context.Context, countryID string, country *model.Country) error {
	return uc.repo.Set(ctx, countryID, country)
}

func (uc *CountryUsecase) ListFromRepo(ctx context.Context) ([]*model.Country, error) {
	keys, err := uc.repo.Keys(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	countries := make([]*model.Country, 0, len(keys))
	for _, key := range keys {
		country, err := uc.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}
	return countries, nil
}

func (uc *CountryUsecase) ListFromCache(ctx context.Context) []*model.Country {
	return uc.repo.Cache(ctx)
}

func (uc *CountryUsecase) GetInfo(ctx context.Context, language language.Tag, country *model.Country) (string, error) {
	return uc.service.CountryInfo(ctx, language, country), nil
}

func (uc *CountryUsecase) DefaultCountry(ctx context.Context) *model.Country {
	return uc.service.DefaultCountry(ctx)
}

func (uc *CountryUsecase) loadFromFiles(ctx context.Context, dir string) error {
	entries, err := assets.Countries.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("unable to read %s dir: %w", dir, err)
	}
	if len(entries) == 0 {
		return country.ErrNoFiles
	}
	for _, f := range entries {
		content, err := fs.ReadFile(assets.Countries, filepath.Join(dir, f.Name()))
		if err != nil {
			return fmt.Errorf("unable to read %s: %w", f.Name(), err)
		}
		var country model.Country
		if err := toml.Unmarshal(content, &country); err != nil {
			return fmt.Errorf("unable to unmarshal %s: %w", f.Name(), err)
		}
		slog.Debug("init for repo", "repo", uc.repo)
		if err := uc.repo.Set(ctx, country.Code, &country); err != nil {
			return err
		}
	}
	return nil
}
