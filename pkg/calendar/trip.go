package calendar

import (
	"fmt"
	"time"

	"cdtj.io/days-in-turkey-bot/user"
)

type TripTree struct {
	StartDate    time.Time
	EndDate      time.Time
	EndPredicted bool
	TripDays     int
	PeriodDays   int

	Prev *TripTree
	Next *TripTree
}

func Trip(u *user.UserConfig, dates []time.Time) {
	tree := calcTree(u, dates)
	for i := tree; i != nil; i = i.Prev {
		fmt.Printf("trip: %q - %q @ %d / %d\n", i.StartDate, i.EndDate, i.TripDays, i.PeriodDays)
	}
}

func calcTree(u *user.UserConfig, dates []time.Time) *TripTree {
	var prev, tree *TripTree
	for i := 0; i < len(dates); i = i + 2 {
		var endDate time.Time
		startDate := dates[i]
		predicted := false
		if i < len(dates)-1 {
			endDate = dates[i+1]
		} else {
			predicted = true
			endDate = predictDate(startDate, maxTripDate(startDate, u.GetDaysCont()), u.GetDaysLimit(), u.GetDaysCont(), u.GetDaysReset(), tree)
		}
		daysPassed := tripDays(startDate, endDate, u.GetDaysReset(), tree)
		tree = &TripTree{
			StartDate:    startDate,
			EndDate:      endDate,
			TripDays:     daysBetween(startDate, endDate),
			PeriodDays:   daysPassed,
			EndPredicted: predicted,
			Prev:         prev,
		}
		prev = tree
	}
	return tree
}

func daysAllowed(daysPassed, daysLimit int) int {
	return daysLimit - daysPassed
}

func predictDate(startDate, endDate time.Time, daysLimit, daysCont, daysReset int, tree *TripTree) time.Time {
	checker := daysCont
	for {
		daysPassed := tripDays(startDate, maxTripDate(startDate, checker), daysReset, tree)
		daysAllowed := daysAllowed(daysPassed, daysLimit)
		// fmt.Printf("[%d], daysPassed: %d, daysAllowed: %d\n", checker, daysPassed, daysAllowed)
		if checker < 0 {
			return startDate
		}
		if daysAllowed > 0 {
			return startDate.Add(time.Hour * 24 * time.Duration(checker))
		}
		// daysAllowed are negative on un-allowance
		checker = checker + daysAllowed - 1
	}
}

func tripDays(startDate, endDate time.Time, resetInterval int, tree *TripTree) int {
	resetDay := getResetDay(endDate, resetInterval)
	daysPassed := 0
	if startDate.Before(endDate) {
		daysPassed = daysBetween(startDate, endDate)
	}
	fmt.Printf("%d > %q - %q @ %q\n", daysPassed, startDate, endDate, resetDay)
	for i := tree; i != nil; i = i.Prev {
		if i.EndDate.After(resetDay) {
			branchStartDate := i.StartDate
			if branchStartDate.Before(resetDay) {
				branchStartDate = resetDay
			}
			daysPassed += daysBetween(branchStartDate, i.EndDate)
			// fmt.Printf("%d > on Before(%q): %q - %q\n", daysPassed, resetDay, i.StartDate, i.EndDate)
		} else if i.StartDate.After(resetDay) {
			daysPassed += daysBetween(i.StartDate, i.EndDate)
			// fmt.Printf("%d > on After(%q): %q - %q\n", daysPassed, resetDay, i.StartDate, i.EndDate)
		} else {
			// fmt.Printf("%d > on Break: %q - %q\n", daysPassed, i.StartDate, i.EndDate)
			break
		}
	}
	return daysPassed
}

func getResetDay(dt time.Time, resetInterval int) time.Time {
	return dt.Add(-1 * time.Hour * 24 * time.Duration(resetInterval))
}

func daysBetween(from, to time.Time) int {
	return int(to.Sub(from) / (time.Hour * 24))
}

func maxTripDate(startDate time.Time, daysCont int) time.Time {
	return startDate.Add(time.Hour * 24 * time.Duration(daysCont))
}
