package calendar

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/i18n"
	"golang.org/x/text/language"
)

func TestCalendarCalc(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
	i18n, err := i18n.NewI18n("i18n", language.English.String())
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Name  string
		Input string
	}{
		{"Ugly Date", "asd"},
		{"Predict End Date", "11/11/2023"},
		{"Predict Eligible", "11/11/2023 11/12/2023"},
		{"Bunch Of Inputs", "11/01/2023 11/02/2023 11/03/2023 11/04/2023 11/05/2023 11/06/2023 11/07/2023 11/08/2023 11/09/2023 11/10/2023"},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tree, err := MakeTree(tc.Input, 60, 90, 180)
			if err != nil {
				if tc.Name == "Ugly Date" && errors.Is(err, ErrInvalidDate) {
					return
				}
				t.Error(err)
				return
			}
			t.Log(formatter.NewTelegramFormatter(i18n, false).TripTree(i18n.DefaultLang(), tree))
		})
	}
}

func BenchmarkCalendarCalc(b *testing.B) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
	i18n, err := i18n.NewI18n("i18n", language.English.String())
	if err != nil {
		b.Fatal(err)
	}

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
				tree, err := MakeTree(tc.Input, 60, 90, 180)
				if err != nil {
					b.Error(err)
					continue
				}
				formatter.NewTelegramFormatter(i18n, false).TripTree(i18n.DefaultLang(), tree)
			}
		})
	}
}
