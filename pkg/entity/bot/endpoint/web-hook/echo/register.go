package webhook

import (
	"cdtj.io/days-in-turkey-bot/entity/bot"
	"github.com/labstack/echo/v4"
)

func RegisterWebhookEndpointsEcho(router *echo.Echo, uc bot.Usecase) {
	h := NewBotWebhookHandlerEcho(uc)
	router.POST("/telegram-webhook", h.webhook)
}
