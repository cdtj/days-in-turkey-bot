package bot

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
)

type Service interface {
	Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error
	CountryMarkup(ctx context.Context, countries []*model.Country) []*model.TelegramBotCommandRow
	LangMarkup(ctx context.Context) []*model.TelegramBotCommandRow
}
