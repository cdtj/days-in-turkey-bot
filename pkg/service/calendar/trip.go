package calendar

import (
	"log/slog"
	"time"

	"cdtj.io/days-in-turkey-bot/model"
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
			endDate = predictEndDate(startDate, daysLimit, daysCont, resetInterval, tree)
		}
		tree = makeTree(startDate, endDate, resetInterval, false, predicted, tree)
		tree.Prev = prev
		prev = tree
	}
	if len(dates)%2 == 0 {
		predStart := predictStartDate(isPastToday(tree.EndDate), daysLimit, daysCont, resetInterval, tree)
		eligibleFullTree := makeTree(predStart, predictEndDate(predStart, daysLimit, daysCont, resetInterval, tree), resetInterval, true, false, tree)
		eligibleCurrentTree := makeTree(isPastToday(tree.EndDate), predictEndDate(isPastToday(tree.EndDate), daysLimit, daysCont, resetInterval, tree), resetInterval, true, false, tree)

		if eligibleCurrentTree.TripDays > 0 && eligibleCurrentTree.TripDays != eligibleFullTree.TripDays {
			tree = eligibleCurrentTree
			tree.Prev = prev
			prev = tree
		}

		tree = eligibleFullTree
		tree.Prev = prev
		prev = tree
	}
	return tree
}

func makeTree(startDate, endDate time.Time, resetInterval int, startPred, endPred bool, tree *model.TripTree) *model.TripTree {
	daysPassed := tripDays(startDate, endDate, resetInterval, tree)
	return &model.TripTree{
		StartDate:      startDate,
		EndDate:        endDate,
		TripDays:       daysBetween(startDate, endDate),
		PeriodDays:     daysPassed,
		StartPredicted: startPred,
		EndPredicted:   endPred,
	}
}

func daysAllowed(daysPassed, daysLimit int) int {
	return daysLimit - daysPassed
}

func predictEndDate(startDate time.Time, daysLimit, daysCont, resetInterval int, tree *model.TripTree) time.Time {
	checker := daysCont
	for {
		departureDay := maxDepartureDate(startDate, checker)
		daysPassed := tripDays(startDate, departureDay, resetInterval, tree)
		daysAllowed := daysAllowed(daysPassed, daysLimit)
		slog.Debug("loop", "method", "predictEndDate", "startDate", startDate, "departureDay", departureDay, "daysPassed", daysPassed, "daysAllowed", daysAllowed)
		if checker < 0 {
			return startDate
		}
		if daysAllowed >= 0 {
			return departureDay
		}
		// daysAllowed are negative on un-allowance
		checker += daysAllowed
	}
}

func predictStartDate(startDate time.Time, daysLimit, daysCont, resetInterval int, tree *model.TripTree) time.Time {
	for {
		departureDay := maxDepartureDate(startDate, daysCont)
		daysPassed := tripDays(startDate, departureDay, resetInterval, tree)
		daysAllowed := daysAllowed(daysPassed, daysLimit)
		slog.Debug("loop", "method", "predictStartDate", "startDate", startDate, "departureDay", departureDay, "daysPassed", daysPassed, "daysAllowed", daysAllowed)
		if daysLimit >= daysPassed {
			return startDate
		}
		startDate = startDate.Add(time.Hour * 24 * time.Duration(daysPassed-daysLimit))
	}
}

// tripDays returns int number of days counted in TripTree between two dates limited with resetInterval
func tripDays(startDate, endDate time.Time, resetInterval int, tree *model.TripTree) int {
	resetDay := getResetDay(endDate, resetInterval)
	daysPassed := 0
	if startDate.Before(endDate) {
		daysPassed = daysBetween(startDate, endDate)
	}
	slog.Debug("input", "method", "tripDays", "daysPassed", daysPassed, "startDate", startDate, "endDate", endDate, "resetDay", resetDay)
	for i := tree; i != nil; i = i.Prev {
		if i.StartPredicted {
			continue
		}
		slog.Debug("walking tree", "method", "tripDays", "daysPassed", daysPassed, "StartDate", i.StartDate, "EndDate", i.EndDate, "resetDay", resetDay)
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

// daysBetween returns date subtracted by resetInterval days
func getResetDay(dt time.Time, resetInterval int) time.Time {
	return dt.Add(-1 * time.Hour * 24 * time.Duration(resetInterval))
}

// daysBetween returns int number of days between to dates
func daysBetween(from, to time.Time) int {
	return int(to.Sub(from) / (time.Hour * 24))
}

// maxDepartureDate returns date appended by daysCont days
func maxDepartureDate(startDate time.Time, daysCont int) time.Time {
	return startDate.Add(time.Hour * 24 * time.Duration(daysCont))
}
