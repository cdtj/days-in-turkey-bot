package calendar

import (
	"fmt"
	"strings"
	"time"

	"cdtj.io/days-in-turkey-bot/user"
	"golang.org/x/exp/slices"
)

func ProcessInput(input string) ([]time.Time, error) {
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
	slices.SortFunc(dates, func(a, b time.Time) int { return a.Compare(b) })

	return dates, nil
}

func parseDate(dt string) (time.Time, error) {
	return time.Parse(dateLayout(dt), dt)
}

func dateLayout(dt string) string {
	sep := dateSeparator(dt)
	dtArr := strings.Split(dt, sep)
	if len(dtArr) == 3 {
		if len(dt) == 10 {
			return strings.Join([]string{"02", "01", "2006"}, sep)
		} else if len(dt) == 8 {
			if len(dtArr[2]) == 4 {
				return strings.Join([]string{"2", "1", "2006"}, sep)
			}
			return strings.Join([]string{"02", "01", "06"}, sep)
		}
	}
	return ""
}

func dateSeparator(dt string) string {
	if strings.Contains(dt, "/") {
		return "/"
	}
	if strings.Contains(dt, ".") {
		return "."
	}
	return "-"
}

func strToChunks(str string) []string {
	str = strings.ReplaceAll(str, ",", " ")
	str = strings.ReplaceAll(str, "\n", " ")
	result := strings.Split(str, " ")
	return result
}

func daysSpent(u *user.UserConfig, dates []time.Time) {
	var daysIn, daysOut int
	resetDate := time.Now().Add(-1 * time.Hour * time.Duration(u.GetDaysReset()))
	fmt.Println("resetDate:", resetDate)

	for i := 0; i < len(dates); i++ {

	}
	time.Now().Sub(resetDate).Hours()
}
