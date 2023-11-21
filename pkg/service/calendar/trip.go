package calendar

import (
	"log/slog"
	"time"

	"cdtj.io/days-in-turkey-bot/model"
)

// calcTree is the main function that returns Trip Tree
func calcTree(daysCont, daysLimit, resetInterval int, trips []model.Trip) *model.TripTree {
	var prev, tree *model.TripTree

	for _, trip := range trips {
		// predicted flag to mark EndDate was predicted
		predicted := false
		if trip.EndDate == nil {
			predicted = true
			predictedEnd := predictEndDate(*trip.StartDate, daysCont, daysLimit, resetInterval, tree)
			trip.EndDate = &predictedEnd
		}
		tree = makeTree(*trip.StartDate, *trip.EndDate, daysCont, daysLimit, resetInterval, false, predicted, tree)
		tree.Prev = prev
		prev = tree
	}
	// if tree isn't null and current trip has ended, we're trying to predict closest possilbe daysCont trip
	// and the max duration possible with unused days
	if tree != nil && !tree.EndPredicted {
		// calculating prediction StartDate with past trip EndDate or Today
		predStart := predictStartDate(isPastToday(tree.EndDate), daysCont, daysLimit, resetInterval, tree)
		predEnd := predictEndDate(predStart, daysCont, daysLimit, resetInterval, tree)
		eligibleFullTree := makeTree(predStart, predEnd, daysCont, daysLimit, resetInterval, true, false, tree)
		// using trip EndDate or Today to calculate trip with unused days
		predEnd = predictEndDate(isPastToday(tree.EndDate), daysCont, daysLimit, resetInterval, tree)
		eligibleCurrentTree := makeTree(isPastToday(tree.EndDate), predEnd, daysCont, daysLimit, resetInterval, true, false, tree)

		// if we don't have any unused days or we can do daysCont trip
		// just skip the result so don't have a zerotrip/duplicate
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

func makeTree(startDate, endDate time.Time, daysCont, daysLimit, resetInterval int, startPred, endPred bool, tree *model.TripTree) *model.TripTree {
	daysPassed := tripDays(startDate, endDate, resetInterval, tree)
	dBetween := daysBetween(startDate, endDate)
	return &model.TripTree{
		StartDate:      startDate,
		EndDate:        endDate,
		TripDays:       dBetween,
		OverstayDays:   daysOverstayed(startDate, endDate, daysCont, daysLimit, daysPassed, endPred),
		PeriodDays:     daysPassed,
		StartPredicted: startPred,
		EndPredicted:   endPred,
	}
}

func daysAllowed(daysPassed, daysLimit int) int {
	return daysLimit - daysPassed
}

func daysOverstayed(startDate, endDate time.Time, daysCont, daysLimit, daysPeriod int, endPredicted bool) int {
	if endPredicted {
		endDate = isPastToday(endDate)
	}
	dBetween := daysBetween(startDate, endDate)
	overstayPeriod := daysPeriod - daysLimit
	overstayCont := dBetween - daysCont
	overstayDays := 0
	if overstayPeriod > overstayDays {
		overstayDays = overstayPeriod
	}
	if overstayCont > overstayDays {
		overstayDays = overstayCont
	}
	slog.Debug("daysOverstayed", "StartDate", startDate, "EndDate", endDate, "dBetween", dBetween, "overstayPeriod", overstayPeriod, "overstayCont", overstayCont)
	return overstayDays
}

// predictEndDate tries to predict leave date to avoid overstay
// can be used in predictions on ongoing and future trips
func predictEndDate(startDate time.Time, daysCont, daysLimit, resetInterval int, tree *model.TripTree) time.Time {
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

// predictStartDate tries to predict closes possible start date
// to fullyfy daysCont trip
func predictStartDate(startDate time.Time, daysCont, daysLimit, resetInterval int, tree *model.TripTree) time.Time {
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

func datesToTrips(dates []time.Time) []model.Trip {
	trips := make([]model.Trip, 0, len(dates)%2)
	for i := 0; i < len(dates); i = i + 2 {
		if i+1 < len(dates) {
			trips = append(trips, model.Trip{
				StartDate: &dates[i],
				EndDate:   &dates[i+1],
			})
		} else {
			trips = append(trips, model.Trip{
				StartDate: &dates[i],
				EndDate:   nil,
			})
		}
	}
	return trips
}
