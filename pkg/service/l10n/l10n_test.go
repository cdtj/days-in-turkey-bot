package l10n

import (
	"fmt"
	"testing"

	"golang.org/x/text/language"
)

func InitTest(t *testing.T) {
	langCodes := []string{"ru", "en"}
	for _, langCode := range langCodes {
		lang, err := language.Parse(langCode)
		if err != nil {
			fmt.Println("unknown lang: ", err)
			continue
		}
		fmt.Println(Localaze(lang, "HelloMessage"))
		for i := 0; i < 5; i++ {
			fmt.Println(LocalazeWithCount(lang, "DaysLeft", i))
		}
	}
}
