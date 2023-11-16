package tghandler

import (
	"context"

	"cdtj.io/days-in-turkey-bot/entity/bot"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type BotHandler struct {
	usecase bot.Usecasev2
}

func NewBotHandler(usecase bot.Usecasev2) *BotHandler {
	return &BotHandler{
		usecase: usecase,
	}
}

func (h *BotHandler) welcome(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   h.usecase.Welcome(ctx, update.Message.Chat.ID, update.Message.From.LanguageCode).Text,
	})
}

func (h *BotHandler) country(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   h.usecase.Country(ctx, update.Message.Chat.ID).Text,
	})
}

func (h *BotHandler) language(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   h.usecase.Language(ctx, update.Message.Chat.ID).Text,
	})
}

func (h *BotHandler) contribute(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   h.usecase.Contribute(ctx, update.Message.Chat.ID).Text,
	})
}

func (h *BotHandler) trip(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   h.usecase.Trip(ctx, update.Message.Chat.ID).Text,
	})
}
