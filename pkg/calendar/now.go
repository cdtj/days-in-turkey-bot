package calendar

import (
	"sync"
	"time"
)

var istanbul *time.Location

func getLocation() *time.Location {
	sync.OnceFunc(func() {
		loc, err := time.LoadLocation("Asia/Istanbul")
		if err != nil {
			return
		}
		istanbul = loc
	})
	return istanbul
}

func getToday() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, getLocation())
}
