package http

import (
	"cdtj.io/days-in-turkey-bot/entity/country"
	"github.com/labstack/echo/v4"
)

func RegisterHTTPEndpointsEcho(router *echo.Echo, uc country.Usecase) {
	h := NewCountryHttpHandlerEcho(uc)
	router.GET("/country/get/:countryID", h.getCountry)
}
