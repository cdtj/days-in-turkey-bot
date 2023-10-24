package calendar

import (
	"testing"

	"cdtj.io/days-in-turkey-bot/user"
	"golang.org/x/text/language"
)

func TestDirSrvStartTime(t *testing.T) {
	dates, err := ProcessInput(`
	11/05/2023
	10/07/2023
	15/07/2023
	11/08/2023
	05/11/2023
`)
	if err != nil {
		t.Error(err)
		return
	}
	u := user.NewUserConfig(language.Russian, user.NewUserCountry("ru", 60, 90, 180))
	t.Log("dates:")
	for _, dt := range dates {
		t.Logf("\t%s", dt)
	}
	t.Logf("len: %d", len(dates))
	Trip(u, dates)
}
