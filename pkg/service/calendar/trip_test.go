package calendar

import (
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestDaysBetween(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	})))

	testCases := []struct {
		Name     string
		From     time.Time
		To       time.Time
		Expected int
	}{
		{"Same Day", time.Now(), time.Now(), 1},
		{"One Day", time.Now(), time.Now().AddDate(0, 0, 1), 2},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			days := daysBetween(tc.From, tc.To)
			if days != tc.Expected {
				t.Errorf("Expected %d days, got %d", tc.Expected, days)
			}
		})
	}
}
