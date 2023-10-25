package calendar

import (
	"fmt"
	"time"

	"cdtj.io/days-in-turkey-bot/model"
	"github.com/sirupsen/logrus"
)

func calcTree(daysLimit, daysCont, resetInterval int, dates []time.Time) *model.TripTree {
	var prev, tree *model.TripTree
	for i := 0; i < len(dates); i = i + 2 {
		var endDate time.Time
		startDate := dates[i]
		predicted := false
		if i < len(dates)-1 {
			endDate = dates[i+1]
		} else {
			predicted = true
			endDate = predictDate(startDate, maxTripDate(startDate, daysCont), daysLimit, daysCont, resetInterval, tree)
		}
		daysPassed := tripDays(startDate, endDate, resetInterval, tree)
		tree = &model.TripTree{
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

func predictDate(startDate, endDate time.Time, daysLimit, daysCont, resetInterval int, tree *model.TripTree) time.Time {
	checker := daysCont
	for {
		daysPassed := tripDays(startDate, maxTripDate(startDate, checker), resetInterval, tree)
		daysAllowed := daysAllowed(daysPassed, daysLimit)
		logrus.WithFields(logrus.Fields{"method": "predictDate", "daysPassed": daysPassed, "daysAllowed": daysAllowed}).Debug(checker)
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

func tripDays(startDate, endDate time.Time, resetInterval int, tree *model.TripTree) int {
	resetDay := getResetDay(endDate, resetInterval)
	daysPassed := 0
	if startDate.Before(endDate) {
		daysPassed = daysBetween(startDate, endDate)
	}
	fmt.Printf("%d > %q - %q @ %q\n", daysPassed, startDate, endDate, resetDay)
	for i := tree; i != nil; i = i.Prev {
		logrus.WithFields(logrus.Fields{"method": "tripDays", "StartDate": i.StartDate, "EndDate": i.EndDate, "resetDay": resetDay}).Debug(daysPassed)
		if i.EndDate.After(resetDay) {
			branchStartDate := i.StartDate
			if branchStartDate.Before(resetDay) {
				branchStartDate = resetDay
			}
			daysPassed += daysBetween(branchStartDate, i.EndDate)
		} else if i.StartDate.After(resetDay) {
			daysPassed += daysBetween(i.StartDate, i.EndDate)
		} else {
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
