package i18n

import (
	"testing"

	"golang.org/x/text/language"
)

func TestHelloMessage(t *testing.T) {
	testCases := []struct {
		Lang   string
		Result string
	}{
		{"en", "This bot is an open-source project. You can contribute to Localization, Data Accuracy, and Source Code aswell. Details: https://cdtj.io/l/turkey-bot"},
		{"ru", "Этот бот является open-source проектом. Вы можете внести свой вклад в локализацию, точность даных, а так же в исходный код проекта. Больше информации: https://cdtj.io/l/turkey-bot"},
		{"es", "This bot is an open-source project. You can contribute to Localization, Data Accuracy, and Source Code aswell. Details: https://cdtj.io/l/turkey-bot"}, // default language for unlocalized tag
	}

	i18n, err := NewI18n("i18n", language.English.String())
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		t.Run(tc.Lang, func(t *testing.T) {
			msg := i18n.GetLocaleByString(tc.Lang).Message("Contribute")
			if msg != tc.Result {
				t.Errorf("Expected %q but got %q for %q\n", tc.Result, msg, tc.Lang)
			}
		})
	}
}
