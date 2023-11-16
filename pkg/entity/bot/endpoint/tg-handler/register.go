package tghandler

import (
	"cdtj.io/days-in-turkey-bot/entity/bot"
	telegrambot "cdtj.io/days-in-turkey-bot/telegram-bot"
)

func RegisterBotHandlers(b *telegrambot.TelegramBotv2, uc bot.Usecasev2) {
	h := NewBotHandler(uc)
	b.RegisterHandlerExactMessage("/start", h.welcome)
	b.RegisterHandlerExactMessage("/country", h.country)
	b.RegisterHandlerExactMessage("/language", h.language)
	b.RegisterHandlerExactMessage("/contribute", h.contribute)
	b.RegisterHandlerExactMessage("/trip", h.trip)
}
