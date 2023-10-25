package model

import "time"

type TripTree struct {
	StartDate    time.Time
	EndDate      time.Time
	EndPredicted bool
	TripDays     int
	PeriodDays   int

	Prev *TripTree
}
