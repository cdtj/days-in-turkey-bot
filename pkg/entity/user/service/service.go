package service

import (
	"context"

	"cdtj.io/days-in-turkey-bot/entity/user"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/calendar"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/l10n"
	"golang.org/x/text/language"
)

var _ user.Service = NewUserService(nil)

type UserService struct {
	fmtr formatter.Formatter
}

func NewUserService(fmtr formatter.Formatter) *UserService {
	return &UserService{
		fmtr: fmtr,
	}
}

func (s *UserService) UserInfo(ctx context.Context, l *l10n.Locale, u *model.User) string {
	return s.fmtr.User(l, u)
}

func (s *UserService) CalculateTrip(ctx context.Context, l *l10n.Locale, input string, daysLimit, daysCont, resetInterval int) (string, error) {
	tree, err := calendar.MakeTree(input, daysLimit, daysCont, resetInterval)
	if err != nil {
		return "", err
	}
	return s.fmtr.TripTree(l, tree), nil
}

func (s *UserService) LangLookup(ctx context.Context, lang string) (language.Tag, error) {
	return language.Parse(lang)
}
