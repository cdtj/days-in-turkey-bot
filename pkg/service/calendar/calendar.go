package calendar

import "cdtj.io/days-in-turkey-bot/model"

type Calendarer interface {
	MakeTree(input string, daysCont, daysLimit, resetInterval int) (*model.TripTree, error)
}

type Calendar struct{}

func NewCalenadr() *Calendar { return &Calendar{} }

// MakeTree calculates Trip Tree
// errors: ErrInvalidDate, ErrInvalidYear, ErrInvalidMonth, ErrInvalidDay
func MakeTree(input string, daysCont, daysLimit, resetInterval int) (*model.TripTree, error) {
	dates, err := processInput(input)
	if err != nil {
		return nil, err
	}
	return calcTree(daysCont, daysLimit, resetInterval, datesToTrips(dates)), nil
}
