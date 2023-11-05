package country

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"cdtj.io/days-in-turkey-bot/db"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/i18n"

	cr "cdtj.io/days-in-turkey-bot/entity/country/repo"
	cs "cdtj.io/days-in-turkey-bot/entity/country/service"
	cuc "cdtj.io/days-in-turkey-bot/entity/country/usecase"
)

const (
	FromCache = iota
	FromSyncMap
)

func BenchmarkCountry(b *testing.B) {
	i18n.Localization()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

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
			countrySvc := cs.NewCountryService(formatter.NewTelegramFormatter())
			countryUC := cuc.NewCountryUsecase(countryRepo, countrySvc)
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				switch tc.Input {
				case FromCache:
					countryUC.Cache(context.Background())
				case FromSyncMap:
					_, err := countryUC.List(context.Background())
					if err != nil {
						b.Error(err)
					}
				}
			}
		})
	}
}
