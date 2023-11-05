package calendar

import (
	"log/slog"
	"os"
	"testing"

	"cdtj.io/days-in-turkey-bot/service/l10n"
)

func TestProcessInput(t *testing.T) {
	l10n.Localization()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

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
			dates, err := processInput(tc.Input)
			if err != nil {
				t.Error(err)
			}
			t.Log(dates)
		})
	}
}
