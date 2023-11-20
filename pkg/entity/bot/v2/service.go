package bot

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
	tgmodel "github.com/go-telegram/bot/models"
	"golang.org/x/text/language"
)

type Service interface {
	FormatMessage(ctx context.Context, language language.Tag, messageID FmtdMsg) string
	CountryMarkup(ctx context.Context, countries []*model.Country) []*model.TelegramBotCommandRow
	LanguageMarkup(ctx context.Context) []*model.TelegramBotCommandRow
	CommandsToInlineKeboard(ctx context.Context, commands []*model.TelegramBotCommandRow) *tgmodel.InlineKeyboardMarkup
}
