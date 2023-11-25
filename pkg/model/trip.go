package model

import "time"

type TripTree struct {
	StartDate      time.Time
	EndDate        time.Time
	EndPredicted   bool
	StartPredicted bool
	TripDays       int
	PeriodDays     int
	OverstayDays   int

	Prev *TripTree
}

type Trip struct {
	startDate *time.Time
	endDate   *time.Time
}

func NewTrip(startDate, endDate *time.Time) Trip {
	return Trip{
		startDate: startDate,
		endDate:   endDate,
	}
}

func (s *Trip) GetStartDate() *time.Time {
	return s.startDate
}
func (s *Trip) SetStartDate(startDate *time.Time) {
	s.startDate = startDate
}
func (s *Trip) GetEndDate() *time.Time {
	return s.endDate
}
func (s *Trip) SetEndDate(endDate *time.Time) {
	s.endDate = endDate
}
