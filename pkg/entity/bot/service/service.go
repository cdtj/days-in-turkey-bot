package service

import (
	"context"

	"cdtj.io/days-in-turkey-bot/entity/bot"
	"cdtj.io/days-in-turkey-bot/model"
)

type BotService struct {
	service bot.Service
}

func NewBotService(service bot.Service) *BotService {
	return &BotService{
		service: service,
	}
}

func (s *BotService) Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error {
	return s.service.Send(ctx, chatID, text, replyMarkup)
}
