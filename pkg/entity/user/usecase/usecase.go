package usecase

import (
	"context"
	"errors"

	"cdtj.io/days-in-turkey-bot/entity/user"
	"cdtj.io/days-in-turkey-bot/model"
)

type UserUsecase struct {
	repo    user.Repo
	service user.Service
}

func NewUserUsecase(repo user.Repo, service user.Service) *UserUsecase {
	return &UserUsecase{
		repo:    repo,
		service: service,
	}
}

func (uc *UserUsecase) Create(ctx context.Context, userID uint64) error {
	return uc.repo.Save(ctx, userID, model.DefaultUser())
}

func (uc *UserUsecase) Get(ctx context.Context, userID uint64) (*model.User, error) {
	u, err := uc.repo.Load(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrRepoUserNotFound) {
			if err := uc.Create(ctx, userID); err != nil {
				return nil, err
			}
			return uc.Get(ctx, userID)
		}
		return nil, err
	}
	return u, nil
}

func (uc *UserUsecase) Calc(ctx context.Context, userID uint64, input string) (string, error) {
	u, err := uc.Get(ctx, userID)
	if err != nil {
		return "", err
	}
	return uc.service.Calc(ctx, input, u.GetDaysLimit(), u.GetDaysCont(), u.GetResetInterval())
}
