package http

import (
	"cdtj.io/days-in-turkey-bot/entity/user"
	"github.com/labstack/echo/v4"
)

func RegisterHTTPEndpointsEcho(router *echo.Echo, uc user.Usecase) {
	h := NewUserHttpHandlerEcho(uc)
	router.GET("/user/info/:userID", h.info)
	router.POST("/user/calc/:userID", h.calculateTrip)
	router.POST("/user/country/:userID", h.updateCountry)
	router.POST("/user/lang/:userID", h.updateLang)
}
