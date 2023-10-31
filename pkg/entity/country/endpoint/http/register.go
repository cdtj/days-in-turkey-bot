package http

import (
	"cdtj.io/days-in-turkey-bot/entity/country"
	httpserver "cdtj.io/days-in-turkey-bot/http-server"
)

// used for debug purposes
func RegisterHTTPEndpoints(router httpserver.HttpServerRouter, uc country.Usecase) {
	h := NewCountryHttpHandler(uc)
	router.HandleFunc("/country/get", h.getCountry)
}
