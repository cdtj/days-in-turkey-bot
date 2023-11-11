package tests

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"cdtj.io/days-in-turkey-bot/db"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/i18n"
	"golang.org/x/text/language"

	cr "cdtj.io/days-in-turkey-bot/entity/country/repo"
	cs "cdtj.io/days-in-turkey-bot/entity/country/service"
	cuc "cdtj.io/days-in-turkey-bot/entity/country/usecase"
)

const (
	FromCache = iota
	FromSyncMap
)

func BenchmarkCountry(b *testing.B) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))
	i18n, err := i18n.NewI18n("i18n", language.English.String())
	if err != nil {
		b.Fatal(err)
	}

	testCases := []struct {
		Name  string
		Input int
	}{
		{"List From SyncMap Repo", FromSyncMap},
		{"List From Cache Repo", FromCache},
	}

	for _, tc := range testCases {
		b.Run(tc.Name, func(b *testing.B) {
			countryDB := db.NewMapDB()
			countryRepo := cr.NewCountryRepo(countryDB)
			countrySvc := cs.NewCountryService(formatter.NewTelegramFormatter(i18n), model.NewCountry("RU", "RU", 90, 60, 180))
			countryUC := cuc.NewCountryUsecase(countryRepo, countrySvc)
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				switch tc.Input {
				case FromCache:
					countryUC.ListFromCache(context.Background())
				case FromSyncMap:
					_, err := countryUC.ListFromRepo(context.Background())
					if err != nil {
						b.Error(err)
					}
				}
			}
		})
	}
}
