package service

import (
	"context"

	"cdtj.io/days-in-turkey-bot/entity/bot"
	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/i18n"
	"golang.org/x/text/language"
)

type BotSender interface {
	Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error
}

type BotService struct {
	sender BotSender
	frmtr  formatter.Formatter
	i18n   i18n.I18ner
}

var _ = NewBotService(nil, nil, nil)

func NewBotService(service BotSender, frmtr formatter.Formatter, i18n i18n.I18ner) *BotService {
	return &BotService{
		sender: service,
		frmtr:  frmtr,
		i18n:   i18n,
	}
}

func (s *BotService) Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error {
	if s.sender == nil {
		return bot.ErrBotNotImpl
	}
	return s.sender.Send(ctx, chatID, text, replyMarkup)
}

func (s *BotService) CommandsFromCountry(ctx context.Context, countries []*model.Country) []*model.TelegramBotCommandRow {
	commands := make([]*model.TelegramBotCommand, 0, len(countries))
	for _, country := range countries {
		commands = append(commands, model.NewTelegramBotCommand(country.GetFlag()+" "+country.GetName(), "country "+country.GetCode(), model.TelegramBotCommandCallbackExact))
	}
	return []*model.TelegramBotCommandRow{model.NewTelegramBotCommandRow(commands, "")}
}

func (s *BotService) CommandsFromLanguage(ctx context.Context) []*model.TelegramBotCommandRow {
	locales := s.i18n.Locales()
	commands := make([]*model.TelegramBotCommand, 0, len(locales))
	for _, cmd := range locales {
		commands = append(commands, model.NewTelegramBotCommand(cmd.Name, "language "+cmd.Tag.String(), model.TelegramBotCommandCallbackExact))
	}
	return []*model.TelegramBotCommandRow{model.NewTelegramBotCommandRow(commands, "")}
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
