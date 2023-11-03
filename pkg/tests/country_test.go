package country

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"cdtj.io/days-in-turkey-bot/db"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/l10n"

	cr "cdtj.io/days-in-turkey-bot/entity/country/repo"
	cs "cdtj.io/days-in-turkey-bot/entity/country/service"
	cuc "cdtj.io/days-in-turkey-bot/entity/country/usecase"
)

func BenchmarkCountry(b *testing.B) {
	l10n.Localization()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	testCases := []struct {
		Name  string
		Input string
	}{
		{"List From SyncMap Repo", ""},
	}

	b.ResetTimer()
	for _, tc := range testCases {
		b.Run(tc.Name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				// country service
				countryDB := db.NewMapDB()
				countryRepo := cr.NewCountryRepo(countryDB)
				countrySvc := cs.NewCountryService(formatter.NewTelegramFormatter())
				countryUC := cuc.NewCountryUsecase(countryRepo, countrySvc)
				_, err := countryUC.List(context.Background())
				if err != nil {
					b.Error(err)
				}
			}
		})
	}
}
