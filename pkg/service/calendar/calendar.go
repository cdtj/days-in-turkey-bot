package calendar

import "cdtj.io/days-in-turkey-bot/model"

func MakeTree(input string, daysLimit, daysCont, resetInterval int) (*model.TripTree, error) {
	dates, err := processInput(input)
	if err != nil {
		return nil, err
	}
	return calcTree(daysLimit, daysCont, resetInterval, dates), nil
}
