package bot

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
)

type Usecase interface {
	Welcome(ctx context.Context, chatID int64, userID, lang string) error
	UpdateLang(ctx context.Context, chatID int64, userID, lang string) error
	UpdateCountry(ctx context.Context, chatID int64, userID, countryID string) error
	CalculateTrip(ctx context.Context, chatID int64, userID, input string) error

	Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error
}
