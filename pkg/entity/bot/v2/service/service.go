package service

import (
	"context"

	"cdtj.io/days-in-turkey-bot/entity/bot/v2"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/i18n"
	tgmodel "github.com/go-telegram/bot/models"
	"golang.org/x/text/language"
)

type BotService struct {
	frmtr formatter.Formatter
	i18n  i18n.I18ner
}

var _ = NewBotService(nil, nil)

func NewBotService(frmtr formatter.Formatter, i18n i18n.I18ner) *BotService {
	return &BotService{
		frmtr: frmtr,
		i18n:  i18n,
	}
}
func (s *BotService) CountryMarkup(ctx context.Context, countries []*model.Country) []*model.TelegramBotCommandRow {
	commands := make([]*model.TelegramBotCommand, 0, len(countries))
	for _, country := range countries {
		commands = append(commands, model.NewTelegramBotCommand(country.GetFlag()+" "+country.GetName(), "country "+country.GetCode()))
	}
	return []*model.TelegramBotCommandRow{model.NewTelegramBotCommandRow(commands)}
}

func (s *BotService) LanguageMarkup(ctx context.Context) []*model.TelegramBotCommandRow {
	locales := s.i18n.Locales()
	commands := make([]*model.TelegramBotCommand, 0, len(locales))
	for _, cmd := range locales {
		commands = append(commands, model.NewTelegramBotCommand(cmd.Name, "language "+cmd.Tag.String()))
	}
	return []*model.TelegramBotCommandRow{model.NewTelegramBotCommandRow(commands)}
}

func (s *BotService) FormatMessage(ctx context.Context, language language.Tag, messageID bot.FmtdMsg) string {
	switch messageID {
	case bot.FmtdMsgWelcome:
		return s.frmtr.Welcome(language)
	case bot.FmtdMsgTripExplanation:
		return s.frmtr.TripExplanation(language)
	default:
		return s.frmtr.FormatMessage(language, string(messageID))
	}
}

func (s *BotService) FormatError(ctx context.Context, language language.Tag, err error) string {
	return s.frmtr.FormatError(language, err)
}

func (s *BotService) CommandsToInlineKeboard(ctx context.Context, commands []*model.TelegramBotCommandRow) *tgmodel.InlineKeyboardMarkup {
	ikbs := make([]tgmodel.InlineKeyboardButton, 0)
	for _, command := range commands {
		ikrs := make([]tgmodel.InlineKeyboardButton, 0)
		for _, command := range command.Commands {
			ikrs = append(ikrs, tgmodel.InlineKeyboardButton{Text: command.Caption, CallbackData: command.Command})
		}
		ikbs = append(ikbs, ikrs...)
	}
	inlineKeyboard := &tgmodel.InlineKeyboardMarkup{InlineKeyboard: [][]tgmodel.InlineKeyboardButton{ikbs}}
	return inlineKeyboard
}
