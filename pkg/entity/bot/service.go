package bot

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
	"golang.org/x/text/language"
)

type Service interface {
	Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error
	FormatMessage(ctx context.Context, language language.Tag, messageID FmtdMsg) string
	CommandsFromCountry(ctx context.Context, countries []*model.Country) []*model.TelegramBotCommandRow
	CommandsFromLanguage(ctx context.Context) []*model.TelegramBotCommandRow
}
