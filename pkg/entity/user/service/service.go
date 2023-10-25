package service

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/calendar"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"golang.org/x/text/language"
)

type UserService struct {
	fmtr formatter.Formatter
}

func NewUserService(fmtr formatter.Formatter) *UserService {
	return &UserService{
		fmtr: fmtr,
	}
}

func (s *UserService) UserInfo(ctx context.Context, u *model.User) string {
	return s.fmtr.User(u)
}

func (s *UserService) CalculateTrip(ctx context.Context, input string, daysLimit, daysCont, resetInterval int) (string, error) {
	tree, err := calendar.MakeTree(input, daysLimit, daysCont, resetInterval)
	if err != nil {
		return "", err
	}
	return s.fmtr.TripTree(tree), nil
}

func (s *UserService) LangLookup(ctx context.Context, lang string) (language.Tag, error) {
	return language.Parse(lang)
}
