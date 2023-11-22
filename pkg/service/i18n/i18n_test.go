package i18n

import (
	"log/slog"
	"os"
	"testing"

	"cdtj.io/days-in-turkey-bot/model"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestLocales(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
	i18n, err := NewI18n("i18n", language.English.String())
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Lang   string
		Result string
	}{
		{"en", "If you're not addicted to github pull-requests just post your feedback or questions to https://t.me/TurkeyDays telegram chat"},
		{"ru", "Публичный чат для вопросов и предложений: https://t.me/TurkeyDays"},
		{"es", "If you're not addicted to github pull-requests just post your feedback or questions to https://t.me/TurkeyDays telegram chat"}, // default language for unlocalized tag
	}

	for _, tc := range testCases {
		t.Run(tc.Lang, func(t *testing.T) {
			msg := i18n.GetLocaleByString(tc.Lang).Message("Feedback")
			assert.Equal(t, tc.Result, msg)
		})
	}
}

func TestExpandableError(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
	i18n, err := NewI18n("i18n", language.English.String())
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		Lang   string
		Input  *model.LError
		Result string
	}{
		{"en", model.NewLError("ErrorInvalidDatePeriod", map[string]any{
			"PeriodName":  model.LErrorExpandable("DatePeriodDay"),
			"PeriodValue": "asd",
			"DateInput":   "asd/11/11",
		}, nil), "invalid day: asd [asd/11/11]"},
		{"ru", model.NewLError("ErrorInvalidDatePeriod", map[string]any{
			"PeriodName":  model.LErrorExpandable("DatePeriodDay"),
			"PeriodValue": "asd",
			"DateInput":   "asd/11/11",
		}, nil), "неправильный день: asd [asd/11/11]"},
		{"es", model.NewLError("ErrorInvalidDatePeriod", map[string]any{
			"PeriodName":  model.LErrorExpandable("DatePeriodDay"),
			"PeriodValue": "asd",
			"DateInput":   "asd/11/11",
		}, nil), "invalid day: asd [asd/11/11]"}, // default language for unlocalized tag
	}

	for _, tc := range testCases {
		t.Run(tc.Lang, func(t *testing.T) {
			l := i18n.GetLocaleByString(tc.Lang)
			msg := l.Error(tc.Input)
			assert.Equal(t, tc.Result, msg, "lang:"+tc.Lang)
		})
	}
}
