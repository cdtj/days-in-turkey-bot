package calendar

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cdtj.io/days-in-turkey-bot/model"
	"golang.org/x/exp/slices"
)

var (
	nonNumberRegexp *regexp.Regexp
)

var (
	ErrInvalidDate  = errors.New("invalid date")
	ErrInvalidYear  = errors.New("invalid year")
	ErrInvalidMonth = errors.New("invalid month")
	ErrInvalidDay   = errors.New("invalid day")
)

func init() {
	if err := compileRegexp(); err != nil {
		panic(err)
	}
}

func compileRegexp() error {
	pattern := "[^0-9]"

	var err error
	nonNumberRegexp, err = regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("error compiling regex: %w", err)
	}
	return nil
}

func processInput(input string) ([]time.Time, error) {
	strDates := strToChunks(input)
	dates := make([]time.Time, 0, len(strDates))
	for _, strDate := range strDates {
		strDate = strings.TrimSpace(strDate)
		if strDate == "" {
			continue
		}
		dt, err := parseDate(strDate)
		if err != nil {
			return nil, err
		}
		dates = append(dates, dt)
	}
	if len(dates) == 0 {
		return nil, model.NewLError("ErrorInvalidDate", map[string]any{"DateInput": ""}, ErrInvalidDate)
	}
	slices.SortFunc(dates, func(a, b time.Time) int { return a.Compare(b) })

	return dates, nil
}

func parseDate(dt string) (time.Time, error) {
	sep := getSeparator(dt)
	dtArr := strings.Split(dt, sep)
	nullTime := time.Time{}
	if len(dtArr) == 3 {
		year, err := strconv.Atoi(dtArr[2])
		if err != nil {
			return nullTime, model.NewLError("ErrorInvalidDatePeriod",
				map[string]any{"DateInput": dt, "PeriodName": "DatePeriodYear", "PeriodValue": dtArr[2]},
				errors.Join(ErrInvalidYear, err))
		}
		if year < 2000 {
			year += 2000
		}
		month, err := strconv.Atoi(dtArr[1])
		if err != nil {
			return nullTime, model.NewLError("ErrorInvalidDatePeriod",
				map[string]any{"DateInput": dt, "PeriodName": "DatePeriodMonth", "PeriodValue": dtArr[1]},
				errors.Join(ErrInvalidMonth, err))
		}
		day, err := strconv.Atoi(dtArr[0])
		if err != nil {
			return nullTime, model.NewLError("ErrorInvalidDatePeriod",
				map[string]any{"DateInput": dt, "PeriodName": "DatePeriodDay", "PeriodValue": dtArr[0]},
				errors.Join(ErrInvalidDay, err))
		}
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
	}
	return nullTime, model.NewLError("ErrorInvalidDate",
		map[string]any{"DateInput": dt},
		ErrInvalidDate)
}

func getSeparator(str string) string {
	nonNumbers := nonNumberRegexp.FindAllString(str, -1)
	if len(nonNumbers) == 2 && nonNumbers[0] == nonNumbers[1] {
		return nonNumbers[0]
	}
	return "INVALID_SEPARATOR"
}

func strToChunks(str string) []string {
	str = strings.ReplaceAll(str, ",", " ")
	str = strings.ReplaceAll(str, "\n", " ")
	result := strings.Split(str, " ")
	return result
}
