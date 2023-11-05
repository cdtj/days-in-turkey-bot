package service

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/formatter"
	"cdtj.io/days-in-turkey-bot/service/l10n"
)

type BotSender interface {
	Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error
}

type BotService struct {
	sender BotSender
	fmtr   formatter.Formatter
}

func NewBotService(sender BotSender, fmtr formatter.Formatter) *BotService {
	return &BotService{
		sender: sender,
		fmtr:   fmtr,
	}
}

func (s *BotService) Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error {
	return s.sender.Send(ctx, chatID, text, replyMarkup)
}

func (s *BotService) FormatMessage(ctx context.Context, l *l10n.Locale, messageID string) string {
	return s.fmtr.FormatMessage(l, messageID)
}
