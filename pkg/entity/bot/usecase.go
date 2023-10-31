package bot

import "context"

type Usecase interface {
	CalculateTrip(ctx context.Context, userID, input string) (string, error)
	UpdateLang(ctx context.Context, userID, lang string) (string, error)
	UpdateCountry(ctx context.Context, userID, countryID string) (string, error)
	Reply(ctx context.Context, chatID int64, text string) error
}
