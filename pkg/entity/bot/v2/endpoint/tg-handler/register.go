package tghandler

import (
	"cdtj.io/days-in-turkey-bot/entity/bot/v2"
	telegrambot "cdtj.io/days-in-turkey-bot/telegram-bot/v2"
	tgapi "github.com/go-telegram/bot"
)

func BindBotHandlers(uc bot.Usecase) []tgapi.Option {
	h := NewBotHandler(uc)
	return []tgapi.Option{
		telegrambot.BindHandlerExactMessage("/start", h.welcome),
		telegrambot.BindHandlerExactMessage("/country", h.country),
		telegrambot.BindHandlerExactMessage("/language", h.language),
		telegrambot.BindHandlerExactMessage("/contribute", h.contribute),
		telegrambot.BindHandlerExactMessage("/trip", h.trip),

		telegrambot.BindHandlerPrefixCb("country", h.updateCountry),
		telegrambot.BindHandlerPrefixCb("language", h.updateLanguage),

		telegrambot.BindHandlerPrefixMessage("/custom", h.updateCountry),

		telegrambot.BindHandlerDefault(h.defaultMessage),
	}
}
