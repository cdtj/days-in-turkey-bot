package calendar

import (
	"time"
)

var (
// istanbul *time.Location
// once     sync.Once
)

// initially, I decided to calc everyting in Turkish time zone,
// but after realized why do I do that?
func getLocation() *time.Location {
	/*
		once.Do(func() {
			loc, err := time.LoadLocation("Asia/Istanbul")
			if err != nil {
				fmt.Println("err:", err)
				return
			}
			fmt.Println("loc loaded:", loc)
			istanbul = loc
		})
		return istanbul
	*/
	return time.UTC
}

func getToday() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, getLocation())
}
