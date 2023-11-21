package service

import (
	"context"

	"cdtj.io/days-in-turkey-bot/entity/country"
	"cdtj.io/days-in-turkey-bot/entity/user"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/calendar"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/i18n"
	"golang.org/x/text/language"
)

var _ user.Service = NewUserService(nil, nil, nil)

type UserService struct {
	fmtr    formatter.Formatter
	i18n    i18n.I18ner
	country country.Service
}

func NewUserService(fmtr formatter.Formatter, i18n i18n.I18ner, country country.Service) *UserService {
	return &UserService{
		fmtr:    fmtr,
		i18n:    i18n,
		country: country,
	}
}

func (s *UserService) UserInfo(ctx context.Context, language language.Tag, user *model.User) string {
	return s.fmtr.User(language, user)
}

func (s *UserService) CalculateTrip(ctx context.Context, language language.Tag, input string, daysCont, daysLimit, resetInterval int) (string, error) {
	tree, err := calendar.MakeTree(input, daysCont, daysLimit, resetInterval)
	if err != nil {
		return "", err
	}
	return s.fmtr.TripTree(language, tree), nil
}

func (s *UserService) DefaultUser(ctx context.Context, userID int64) *model.User {
	return model.NewUser(userID, s.i18n.DefaultLang().String(), *s.country.DefaultCountry(ctx))
}

func (s *UserService) ParseLanguage(ctx context.Context, languageCode string) language.Tag {
	tag, err := i18n.LanguageLookup(languageCode)
	if err != nil {
		return s.i18n.DefaultLang()
	}
	return tag
}
