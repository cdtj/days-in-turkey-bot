package bot

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/l10n"
)

type Service interface {
	Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error
	FormatMessage(ctx context.Context, l *l10n.Locale, messageID string) string
}
