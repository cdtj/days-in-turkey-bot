package webhook

import (
	"cdtj.io/days-in-turkey-bot/entity/bot"
	httpserver "cdtj.io/days-in-turkey-bot/http-server"
)

// used for debug purposes
func RegisterWebhookEndpoints(router httpserver.HttpServerRouter, uc bot.Usecase) {
	h := NewBotWebhookHandler(uc)
	router.HandleFunc("/telegram-webhook", h.webhook)
}
