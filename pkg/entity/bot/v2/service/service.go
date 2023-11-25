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

const (
	inlineKeyboardCountryMaxLen  = 3
	inlineKeyboardLanguageMaxLen = 3
)

// CommandsFromCountry creates list of Telegram Bot Commands from the slice of Countries,
// used for inlineKeyboard generation
func (s *BotService) CommandsFromCountry(ctx context.Context, countries []*model.Country) []*model.TelegramBotCommandRow {
	commandsRows := make([]*model.TelegramBotCommandRow, 0)
	commands := make([]*model.TelegramBotCommand, 0)
	for k, country := range countries {
		newline := k + 1
		commands = append(commands, model.NewTelegramBotCommand(country.GetFlag()+" "+country.GetName(), "country "+country.GetCode(), model.TelegramBotCommandCallbackExact))
		if (newline > 0 && newline%inlineKeyboardCountryMaxLen == 0) || newline == len(countries) {
			commandsRows = append(commandsRows, model.NewTelegramBotCommandRow(commands, ""))
			commands = make([]*model.TelegramBotCommand, 0)
		}
	}
	return commandsRows
}

// CommandsFromLanguage creates list of Telegram Bot Commands from the nested i18n pkg,
// used for inlineKeyboard generation
func (s *BotService) CommandsFromLanguage(ctx context.Context) []*model.TelegramBotCommandRow {
	locales := s.i18n.Locales()
	commandsRows := make([]*model.TelegramBotCommandRow, 0)
	commands := make([]*model.TelegramBotCommand, 0)
	for k, locale := range locales {
		newline := k + 1
		commands = append(commands, model.NewTelegramBotCommand(locale.Name, "language "+locale.Tag.String(), model.TelegramBotCommandCallbackExact))
		if (newline > 0 && newline%inlineKeyboardLanguageMaxLen == 0) || newline == len(locales) {
			commandsRows = append(commandsRows, model.NewTelegramBotCommandRow(commands, ""))
			commands = make([]*model.TelegramBotCommand, 0)
		}
	}
	return commandsRows
}

// FormatMessage formats messages to a specific format depending on nested formatter,
// we have some enums here and default method
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

// FormatError formats error to a specific format depending on nested formatter
func (s *BotService) FormatError(ctx context.Context, language language.Tag, err error) string {
	return s.frmtr.FormatError(language, err)
}

// CommandsToInlineKeboard converts Telegram Bot Commands into Telegram Bot API InlineKeyboardMarkup
func (s *BotService) CommandsToInlineKeboard(ctx context.Context, commands []*model.TelegramBotCommandRow) *tgmodel.InlineKeyboardMarkup {
	ikbs := make([][]tgmodel.InlineKeyboardButton, 0)
	for _, command := range commands {
		ikrs := make([]tgmodel.InlineKeyboardButton, 0)
		for _, command := range command.Commands {
			ikrs = append(ikrs, tgmodel.InlineKeyboardButton{Text: command.Caption, CallbackData: command.Command})
		}
		ikbs = append(ikbs, ikrs)
	}
	inlineKeyboard := &tgmodel.InlineKeyboardMarkup{InlineKeyboard: ikbs}
	return inlineKeyboard
}
