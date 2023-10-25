package usecase

import (
	"context"
	"errors"

	"cdtj.io/days-in-turkey-bot/entity/country"
	"cdtj.io/days-in-turkey-bot/entity/user"
	"cdtj.io/days-in-turkey-bot/model"
)

type UserUsecase struct {
	repo    user.Repo
	country country.Repo
	service user.Service
}

func NewUserUsecase(repo user.Repo, country country.Repo, service user.Service) *UserUsecase {
	return &UserUsecase{
		repo:    repo,
		country: country,
		service: service,
	}
}

func (uc *UserUsecase) Get(ctx context.Context, userID string) (string, error) {
	u, err := uc.get(ctx, userID)
	if err != nil {
		return "", err
	}
	return uc.service.UserInfo(ctx, u), nil
}

func (uc *UserUsecase) UpdateLang(ctx context.Context, userID string, lang string) error {
	u, err := uc.get(ctx, userID)
	if err != nil {
		return err
	}
	langTag, err := uc.service.LangLookup(ctx, lang)
	if err != nil {
		return err
	}
	u.Lang = langTag
	return uc.repo.Save(ctx, userID, u)
}

func (uc *UserUsecase) UpdateCountry(ctx context.Context, userID string, countryID string) error {
	u, err := uc.get(ctx, userID)
	if err != nil {
		return err
	}
	country, err := uc.country.Load(ctx, countryID)
	if err != nil {
		return err
	}
	u.Country = country
	return uc.repo.Save(ctx, userID, u)
}

func (uc *UserUsecase) Calc(ctx context.Context, userID string, input string) (string, error) {
	u, err := uc.get(ctx, userID)
	if err != nil {
		return "", err
	}
	return uc.service.CalculateTrip(ctx, input, u.GetDaysLimit(), u.GetDaysCont(), u.GetResetInterval())
}

func (uc *UserUsecase) сreate(ctx context.Context, userID string) error {
	return uc.repo.Save(ctx, userID, model.DefaultUser())
}

func (uc *UserUsecase) get(ctx context.Context, userID string) (*model.User, error) {
	u, err := uc.repo.Load(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrRepoUserNotFound) {
			if err := uc.сreate(ctx, userID); err != nil {
				return nil, err
			}
			return uc.get(ctx, userID)
		}
		return nil, err
	}
	return u, nil
}
