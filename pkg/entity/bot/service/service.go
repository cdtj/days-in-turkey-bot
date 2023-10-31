package service

import (
	"context"

	"cdtj.io/days-in-turkey-bot/entity/bot"
)

type BotService struct {
	sender bot.ServiceSender
}

func NewBotService(sender bot.ServiceSender) *BotService {
	return &BotService{
		sender: sender,
	}
}

func (s *BotService) Reply(ctx context.Context, chatID int64, text string) error {
	return s.sender.Send(ctx, chatID, text)
}
