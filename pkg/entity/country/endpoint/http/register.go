package http

import (
	"cdtj.io/days-in-turkey-bot/entity/country"
	"github.com/go-chi/chi/v5"
	"github.com/labstack/echo/v4"
)

func RegisterHTTPEndpointsEcho(router *echo.Echo, uc country.Usecase) {
	h := NewCountryHttpHandlerEcho(uc)
	router.GET("/country/get/:countryID", h.getCountry)
}

// RegisterHTTPEndpointsChi deprecated cuz echo is more handy in my scenario
func RegisterHTTPEndpointsChi(router *chi.Mux, uc country.Usecase) {
	h := NewCountryHttpHandlerChi(uc)
	router.HandleFunc("/country/get/{countryID}", h.getCountry)
}
