package bot

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
)

type Usecase interface {
	Welcome(ctx context.Context, chatID int64, userID int64, language string) error
	Prompt(ctx context.Context, chatID int64, userID int64, prompt string) error
	UpdateLanguage(ctx context.Context, chatID int64, userID int64, language string) error
	UpdateCountry(ctx context.Context, chatID int64, userID int64, countryID string, daysCont, daysLimit, resetInterval int) error
	CalculateTrip(ctx context.Context, chatID int64, userID int64, datesInput string) error

	Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error
}
