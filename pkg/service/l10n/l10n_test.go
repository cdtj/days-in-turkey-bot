package l10n

import (
	"testing"

	"golang.org/x/text/language"
)

func TestHelloMessage(t *testing.T) {
	if err := Localization(); err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		Lang   string
		Result string
	}{
		{"en", "Hello World!"},
		{"ru", "Привет Мир!"},
		{"es", "Hello World!"}, // default language for unlocalized tag
	}

	for _, tc := range testCases {
		t.Run(tc.Lang, func(t *testing.T) {
			tag, err := language.Parse(tc.Lang)
			if err != nil {
				t.Error("invalid tag", err)
				return
			}
			msg := GetLocale(tag).Message("HelloMessage")
			if msg != tc.Result {
				t.Errorf("Expected %q but got %q for %q\n", tc.Result, msg, tc.Lang)
			}
		})
	}
}
