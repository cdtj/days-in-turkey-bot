package service

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/i18n"
)

type BotSender interface {
	Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error
}

type BotService struct {
	sender BotSender
	frmtr  formatter.Formatter
	i18n   i18n.Localizer
}

func NewBotService(service BotSender) *BotService {
	return &BotService{
		sender: service,
	}
}

func (s *BotService) Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error {
	return s.sender.Send(ctx, chatID, text, replyMarkup)
}

func (s *BotService) CountryMarkup(ctx context.Context, countries []*model.Country) []*model.TelegramBotCommandRow {
	commands := make([]*model.TelegramBotCommand, 0, len(countries))
	for _, country := range countries {
		commands = append(commands, model.NewTelegramBotCommand(country.GetFlag()+" "+country.GetName(), "country "+country.GetCode()))
	}
	return []*model.TelegramBotCommandRow{model.NewTelegramBotCommandRow(commands)}
}

func (s *BotService) LangMarkup(ctx context.Context) []*model.TelegramBotCommandRow {
	locales := s.i18n.Locales()
	commands := make([]*model.TelegramBotCommand, 0, len(locales))
	for _, cmd := range locales {
		commands = append(commands, model.NewTelegramBotCommand(cmd.Name, "language "+cmd.Tag.String()))
	}
	return []*model.TelegramBotCommandRow{model.NewTelegramBotCommandRow(commands)}
}
