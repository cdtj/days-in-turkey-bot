package usecase

import (
	"context"
	"errors"
	"log/slog"

	"cdtj.io/days-in-turkey-bot/entity/country"
	"cdtj.io/days-in-turkey-bot/entity/user"
	"cdtj.io/days-in-turkey-bot/model"
	"golang.org/x/text/language"
)

var _ user.Usecase = NewUserUsecase(nil, nil, nil)

type UserUsecase struct {
	repo      user.Repo
	service   user.Service
	countryUC country.Usecase
}

func NewUserUsecase(repo user.Repo, service user.Service, countryUC country.Usecase) *UserUsecase {
	return &UserUsecase{
		repo:      repo,
		service:   service,
		countryUC: countryUC,
	}
}

func (uc *UserUsecase) Create(ctx context.Context, userID int64, lang string) error {
	return uc.сreate(ctx, userID, lang)
}

func (uc *UserUsecase) Get(ctx context.Context, userID int64) (*model.User, error) {
	return uc.get(ctx, userID)
}

func (uc *UserUsecase) GetInfo(ctx context.Context, user *model.User) string {
	return uc.service.UserInfo(ctx, user.GetLanguage(), user)
}

func (uc *UserUsecase) CalculateTrip(ctx context.Context, user *model.User, datesInput string) (string, error) {
	mth := "CalculateTrip"
	trip, err := uc.service.CalculateTrip(ctx, user.GetLanguage(), datesInput, user.GetDaysCont(), user.GetDaysLimit(), user.GetResetInterval())
	if err != nil {
		slog.Error("usecase failed", "method", mth, "user", user, "datesInput", datesInput, "err", err)
		return "", err
	}
	return trip, nil
}

func (uc *UserUsecase) GetLanguage(ctx context.Context, user *model.User) language.Tag {
	return user.GetLanguage()
}

func (uc *UserUsecase) UpdateLanguage(ctx context.Context, user *model.User, languageCode string) error {
	mth := "UpdateLanguage"
	userID := user.GetID()
	slog.Debug("updateLang", "userID", userID, "languageCode", languageCode)
	user.SetLanguageCode(languageCode)
	user.SetLanguage(uc.service.ParseLanguage(ctx, user.GetLanguageCode()))
	if err := uc.repo.Save(ctx, userID, user); err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "err", err)
		return err
	}
	return nil
}

func (uc *UserUsecase) UpdateCountry(ctx context.Context, user *model.User, country *model.Country) error {
	mth := "UpdateCountry"
	userID := user.GetID()
	slog.Debug("updateCountry", "userID", userID, "country", country)
	user.SetCountry(*country)
	if err := uc.repo.Save(ctx, userID, user); err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "err", err)
		return err
	}
	return nil
}

func (uc *UserUsecase) сreate(ctx context.Context, userID int64, languageCode string) (err error) {
	u := uc.service.DefaultUser(ctx, userID)
	u.SetLanguageCode(languageCode)
	// we don't need to perform this check since the default locale will be picked on userRepo load,
	// so, someday users with previously unsupported locales will see localized messages once they are translated
	/*
		if languageCode != "" {
			tag, err := i18n.LanguageLookup(languageCode)
			if err != nil {
				slog.Error("falied to init user with custom lang", "userID", userID, "languageCode", languageCode, "err", err)
			} else {
				u.SetLanguageCode(languageCode)
				slog.Debug("user uc", "userID", userID, "languageCode", languageCode, "tag", tag)
			}
		}
	*/
	mth := "create"
	if err := uc.repo.Save(ctx, userID, u); err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "err", err)
		return err
	}
	return nil
}

// get is the most expensive method because it contains user constructor performing
// CountryLookup and LanguageLookup, it might be cheaper just to store them instead of
// looking up every time
func (uc *UserUsecase) get(ctx context.Context, userID int64) (*model.User, error) {
	mth := "get"
	u, err := uc.repo.Load(ctx, userID)
	if err != nil {
		// slog.Error("user uc", "userID", userID, "err", err)
		if errors.Is(err, user.ErrRepoUserNotFound) {
			if err := uc.сreate(ctx, userID, ""); err != nil {
				return nil, err
			}
			return uc.get(ctx, userID)
		}
		slog.Error("usecase failed", "method", mth, "userID", userID, "err", err)
		return nil, err
	}

	country, err := uc.countryUC.Get(ctx, u.Country.Code)
	if err != nil {
		slog.Error("usecase failed", "method", mth, "userID", userID, "countryCode", u.Country.Code, "err", err)
		return nil, err
	}
	if country != nil {
		u.SetCountry(*country)
	}
	u.SetLanguage(uc.service.ParseLanguage(ctx, u.GetLanguageCode()))

	return u, nil
}
