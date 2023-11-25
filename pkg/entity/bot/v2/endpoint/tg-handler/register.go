package tghandler

import (
	"cdtj.io/days-in-turkey-bot/entity/bot/v2"
	"cdtj.io/days-in-turkey-bot/model"
	telegrambot "cdtj.io/days-in-turkey-bot/telegram-bot/v2"
)

// BindBotHandlers binds bot.Usecase with Telegram Bot Commands,
// actually this is the entry point where Telegram Bot Commands are defined
func BindBotHandlers(uc bot.Usecase) []*telegrambot.TelegramBotBind {
	h := NewBotHandler(uc)
	return []*telegrambot.TelegramBotBind{
		telegrambot.NewTelegramBotBind(model.NewTelegramBotCommand("", "start", model.TelegramBotCommandMessageExact),
			h.welcome),
		telegrambot.NewTelegramBotBind(model.NewTelegramBotCommand("CommandMe", "me", model.TelegramBotCommandMessageExact),
			h.me),
		telegrambot.NewTelegramBotBind(model.NewTelegramBotCommand("CommandCountry", "country", model.TelegramBotCommandMessageExact),
			h.country),
		telegrambot.NewTelegramBotBind(model.NewTelegramBotCommand("CommandLanguage", "language", model.TelegramBotCommandMessageExact),
			h.language),
		telegrambot.NewTelegramBotBind(model.NewTelegramBotCommand("CommandTrip", "trip", model.TelegramBotCommandMessageExact),
			h.trip),
		telegrambot.NewTelegramBotBind(model.NewTelegramBotCommand("CommandContribute", "contribute", model.TelegramBotCommandMessageExact),
			h.contribute),
		telegrambot.NewTelegramBotBind(model.NewTelegramBotCommand("CommandFeedback", "feedback", model.TelegramBotCommandMessageExact),
			h.feedback),

		telegrambot.NewTelegramBotBind(model.NewTelegramBotCommand("", "/custom", model.TelegramBotCommandMessagePrefix),
			h.updateCountry),

		telegrambot.NewTelegramBotBind(model.NewTelegramBotCommand("", "country", model.TelegramBotCommandCallbackExact),
			h.updateCountry),
		telegrambot.NewTelegramBotBind(model.NewTelegramBotCommand("", "language", model.TelegramBotCommandCallbackExact),
			h.updateLanguage),

		telegrambot.NewTelegramBotBind(model.NewTelegramBotCommand("", "", model.TelegramBotCommandDefaultHandler),
			h.defaultMessage),
	}
}
