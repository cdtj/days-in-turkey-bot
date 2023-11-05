package bot

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/i18n"
)

type Service interface {
	Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error
	FormatMessage(ctx context.Context, l *i18n.Locale, messageID string) string
	CountryMarkup(ctx context.Context, countries []*model.Country) []*model.TelegramBotCommandRow
	LangMarkup(ctx context.Context) []*model.TelegramBotCommandRow
}
