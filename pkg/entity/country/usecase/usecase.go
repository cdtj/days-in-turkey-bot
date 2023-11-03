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
	"cdtj.io/days-in-turkey-bot/service/l10n"
	"github.com/BurntSushi/toml"
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
	if repo != nil {
		if err := uc.InitData(context.Background(), countryFiles); err != nil {
			slog.Error("failed to init CountryUsecase", "err", err)
			return nil
		}
	}
	return uc
}

func (u *CountryUsecase) Get(ctx context.Context, countryID string) (*model.Country, error) {
	return u.repo.Get(ctx, countryID)
}

func (u *CountryUsecase) Set(ctx context.Context, countryID string, country *model.Country) error {
	return u.repo.Set(ctx, countryID, country)
}

func (u *CountryUsecase) Keys(ctx context.Context) ([]string, error) {
	return u.repo.Keys(ctx)
}

func (u *CountryUsecase) List(ctx context.Context) ([]*model.Country, error) {
	keys, err := u.repo.Keys(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	countries := make([]*model.Country, 0, len(keys))
	for _, key := range keys {
		country, err := u.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}
	return countries, nil
}

func (u *CountryUsecase) Info(ctx context.Context, countryID string) (string, error) {
	c, err := u.Get(ctx, countryID)
	if err != nil {
		return "", err
	}
	// default locale is enough for debugging
	return u.service.Info(ctx, l10n.GetLocale(l10n.DefaultLang()), c), nil
}

func (u *CountryUsecase) InitData(ctx context.Context, dir string) error {
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
		slog.Debug("init for repo", "repo", u.repo)
		if err := u.repo.Set(ctx, country.Code, &country); err != nil {
			return err
		}
	}
	return nil
}
