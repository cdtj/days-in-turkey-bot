package webhook

import (
	"cdtj.io/days-in-turkey-bot/entity/bot"
	"github.com/go-chi/chi/v5"
)

// RegisterWebhookEndpointsChi deprecated cuz echo is more handy in my scenario
func RegisterWebhookEndpointsChi(router *chi.Mux, uc bot.Usecase) {
	h := NewBotWebhookHandlerChi(uc)
	router.HandleFunc("/telegram-webhook", h.webhook)
}
