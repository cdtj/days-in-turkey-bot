package webhook

import (
	"cdtj.io/days-in-turkey-bot/entity/bot"
	"github.com/go-chi/chi/v5"
)

// used for debug purposes
func RegisterWebhookEndpoints(router *chi.Mux, uc bot.Usecase) {
	h := NewBotWebhookHandler(uc)
	router.HandleFunc("/telegram-webhook", h.webhook)
}
