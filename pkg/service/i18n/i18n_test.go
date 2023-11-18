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
		{"en", "This bot is an open-source project. You can contribute to Localization, Data Accuracy, and Source Code as well. Details: https://cdtj.io/l/turkey-bot"},
		{"ru", "Этот бот является open-source проектом. Вы можете внести свой вклад в локализацию, точность данных, а так же в исходный код проекта. Больше информации: https://cdtj.io/l/turkey-bot"},
		{"es", "This bot is an open-source project. You can contribute to Localization, Data Accuracy, and Source Code as well. Details: https://cdtj.io/l/turkey-bot"}, // default language for unlocalized tag
	}

	for _, tc := range testCases {
		t.Run(tc.Lang, func(t *testing.T) {
			msg := i18n.GetLocaleByString(tc.Lang).Message("Contribute")
			assert.Equal(t, msg, tc.Result)
		})
	}
}
