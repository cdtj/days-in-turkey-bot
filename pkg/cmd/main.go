package main

import (
	"fmt"

	"cdtj.io/days-in-turkey-bot/service/l10n"
	"golang.org/x/text/language"
)

func init() {
	if err := l10n.Localization(); err != nil {
		panic(err)
	}
}

func main() {
	langCodes := []string{"ru", "en"}
	for _, langCode := range langCodes {
		lang, err := language.Parse(langCode)
		if err != nil {
			fmt.Println("unknown lang: ", err)
			continue
		}
		fmt.Println(l10n.Localaze(lang, "HelloMessage"))
		for i := 0; i < 5; i++ {
			fmt.Println(l10n.LocalazeWithCount(lang, "DaysLeft", i))
		}
	}
}
