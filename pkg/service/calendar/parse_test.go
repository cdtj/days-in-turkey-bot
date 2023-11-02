package calendar

import (
	"log/slog"
	"os"
	"testing"

	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/l10n"
)

func TestCalendarCalc(t *testing.T) {
	l10n.Localization()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	testCases := []struct {
		Name  string
		Input string
	}{
		{"Predict End Date", "11/11/2023"},
		{"Predict Eligible", "11/11/2023 11/12/2023"},
		{"Bunch Of Inputs", "11/01/2023 11/02/2023 11/03/2023 11/04/2023 11/05/2023 11/06/2023 11/07/2023 11/08/2023 11/09/2023 11/10/2023"},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tree, err := MakeTree(tc.Input, 90, 60, 180)
			if err != nil {
				t.Error(err)
			}
			formatter.NewTelegramFormatter().TripTree(l10n.NewLocale(l10n.DefaultLang()), tree)
		})
	}
}

func BenchmarkCalendarCalc(b *testing.B) {
	l10n.Localization()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	testCases := []struct {
		Name  string
		Input string
	}{
		{"Predict End Date", "11/11/2023"},
		{"Predict Eligible", "11/11/2023 11/12/2023"},
		{"Bunch Of Inputs", "11/01/2023 11/02/2023 11/03/2023 11/04/2023 11/05/2023 11/06/2023 11/07/2023 11/08/2023 11/09/2023 11/10/2023"},
	}

	b.ResetTimer()
	for _, tc := range testCases {
		b.Run(tc.Name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				tree, err := MakeTree(tc.Input, 90, 60, 180)
				if err != nil {
					b.Error(err)
				}
				formatter.NewTelegramFormatter().TripTree(l10n.NewLocale(l10n.DefaultLang()), tree)
			}
		})
	}
}
