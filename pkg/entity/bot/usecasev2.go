package bot

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
)

type Usecasev2 interface {
	Welcome(ctx context.Context, userID int64, lang string) *model.TelegramMessage
	Country(ctx context.Context, userID int64) *model.TelegramMessage
	Language(ctx context.Context, userID int64) *model.TelegramMessage
	Contribute(ctx context.Context, userID int64) *model.TelegramMessage
	Trip(ctx context.Context, userID int64) *model.TelegramMessage
	UpdateLanguage(ctx context.Context, userID int64, languageCode string) *model.TelegramMessage
	UpdateCountry(ctx context.Context, userID int64, countryID string, daysCont, daysLimit, resetInterval int) *model.TelegramMessage
	CalculateTrip(ctx context.Context, chatID int64, userID int64, datesInput string) *model.TelegramMessage
}
