package i18n

import (
	"log/slog"
	"os"
	"testing"

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
		{"en", "If you're not addicted to github.com just post your feedback or questions to https://t.me/TurkeyDays telegram chat"},
		{"ru", "Публичный чат для вопросов и предложений: https://t.me/TurkeyDays"},
		{"es", "If you're not addicted to github.com just post your feedback or questions to https://t.me/TurkeyDays telegram chat"}, // default language for unlocalized tag
	}

	for _, tc := range testCases {
		t.Run(tc.Lang, func(t *testing.T) {
			msg := i18n.GetLocaleByString(tc.Lang).Message("Feedback")
			assert.Equal(t, msg, tc.Result)
		})
	}
}
