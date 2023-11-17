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
	StartDate *time.Time
	EndDate   *time.Time
}
