package webhook

import (
	"cdtj.io/days-in-turkey-bot/entity/bot"
	"github.com/go-chi/chi/v5"
	"github.com/labstack/echo/v4"
)

func RegisterWebhookEndpointsEcho(router *echo.Echo, uc bot.Usecase) {
	h := NewBotWebhookHandlerEcho(uc)
	router.POST("/telegram-webhook", h.webhook)
}

// RegisterWebhookEndpointsChi deprecated cuz echo is more handy in my scenario
func RegisterWebhookEndpointsChi(router *chi.Mux, uc bot.Usecase) {
	h := NewBotWebhookHandlerChi(uc)
	router.HandleFunc("/telegram-webhook", h.webhook)
}
