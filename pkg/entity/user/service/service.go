package service

import (
	"context"
	"fmt"

	"cdtj.io/days-in-turkey-bot/service/calendar"
)

type UserService struct {
}

func (uc *UserService) Calc(ctx context.Context, input string, daysLimit, daysCont, resetInterval int) (string, error) {
	dates, err := calendar.ProcessInput(input)
	if err != nil {
		return "", err
	}
	tree := calendar.Trip(daysLimit, daysCont, resetInterval, dates)
	result := ""
	for i := tree; i != nil; i = i.Prev {
		result += fmt.Sprintf("trip: %q - %q @ %d / %d\n", i.StartDate, i.EndDate, i.TripDays, i.PeriodDays)
	}
	return result, nil
}
