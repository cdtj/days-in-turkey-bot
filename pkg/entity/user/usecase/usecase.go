package usecase

import (
	"context"
	"errors"
	"log/slog"

	"cdtj.io/days-in-turkey-bot/entity/country"
	"cdtj.io/days-in-turkey-bot/entity/user"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/l10n"
)

var _ user.Usecase = NewUserUsecase(nil, nil, nil)

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

func (uc *UserUsecase) Create(ctx context.Context, userID, lang string) error {
	if _, err := uc.get(ctx, userID); err == nil {
		return nil
	}
	return uc.сreate(ctx, userID, lang)
}

func (uc *UserUsecase) Info(ctx context.Context, userID string) (string, error) {
	u, err := uc.get(ctx, userID)
	if err != nil {
		slog.Error("get user failed", "userID", userID, "err", err)
		return "", err
	}
	slog.Info("get user", "userid", userID, "data", u)
	return uc.service.UserInfo(ctx, l10n.GetLocale(u.GetLangTag()), u), nil
}

func (uc *UserUsecase) UpdateLang(ctx context.Context, userID string, lang string) error {
	u, err := uc.get(ctx, userID)
	if err != nil {
		slog.Error("updateLang failed", "userID", userID, "lang", lang, "err", err)
		return err
	}
	langTag, err := uc.service.LangLookup(ctx, lang)
	if err != nil {
		slog.Error("updateLang failed", "userID", userID, "lang", lang, "err", err)
		return err
	}
	u.Lang = langTag.String()
	slog.Info("update lang", "userid", userID, "lang", langTag)
	if err := uc.repo.Save(ctx, userID, u); err != nil {
		return err
	}
	return nil
}

func (uc *UserUsecase) UpdateCountry(ctx context.Context, userID string, countryID string) error {
	u, err := uc.get(ctx, userID)
	if err != nil {
		slog.Error("updateCountry failed", "userID", userID, "countryID", countryID, "err", err)
		return err
	}
	country, err := uc.country.Get(ctx, countryID)
	if err != nil {
		slog.Error("updateCountry failed", "userID", userID, "countryID", countryID, "err", err)
		return err
	}
	slog.Info("update country", "userid", userID, "country", country)
	u.Country = *country
	if err := uc.repo.Save(ctx, userID, u); err != nil {
		return err
	}
	return nil
}

func (uc *UserUsecase) CalculateTrip(ctx context.Context, userID string, input string) (string, error) {
	u, err := uc.get(ctx, userID)
	if err != nil {
		slog.Error("calculate trip failed", "userID", userID, "input", input, "err", err)
		return "", err
	}
	return uc.service.CalculateTrip(ctx, l10n.GetLocale(u.GetLangTag()), input, u.GetDaysLimit(), u.GetDaysCont(), u.GetResetInterval())
}

func (uc *UserUsecase) сreate(ctx context.Context, userID, lang string) error {
	u := model.DefaultUser()
	if lang != "" {
		tag, err := uc.service.LangLookup(ctx, lang)
		if err != nil {
			slog.Error("falied to init user with custom lang", "userID", userID, "lang", lang, "err", err)
		} else {
			u.Lang = tag.String()
		}
	}
	return uc.repo.Save(ctx, userID, u)
}

// get is a most expensive method because it contains user constructor performing
// CountryLookup and LanguageLookup, it might be cheaper just to store them instead of
// looking up every time
func (uc *UserUsecase) get(ctx context.Context, userID string) (*model.User, error) {
	u, err := uc.repo.Load(ctx, userID)
	if err != nil {
		if errors.Is(err, user.ErrRepoUserNotFound) {
			if err := uc.сreate(ctx, userID, ""); err != nil {
				return nil, err
			}
			return uc.get(ctx, userID)
		}
		slog.Error("internal user get failed", "userID", userID, "err", err)
		return nil, err
	}

	tag, err := uc.service.LangLookup(ctx, u.Lang)
	if err != nil {
		tag = l10n.DefaultLang()
	}
	u.SetLangTag(tag)

	switch u.Country.Code {
	case "CUSTOM":
		break
	case "":
		u.Country = *model.DefaultCountry()
	default:
		country, err := uc.country.Get(ctx, u.Country.Code)
		if err != nil {
			return nil, err
		}
		u.Country = *country
	}

	return u, nil
}
