package http

import (
	"cdtj.io/days-in-turkey-bot/entity/country"
	"github.com/go-chi/chi/v5"
)

// used for debug purposes
func RegisterHTTPEndpoints(router *chi.Mux, uc country.Usecase) {
	h := NewCountryHttpHandler(uc)
	router.HandleFunc("/country/get/{countryID}", h.getCountry)
}
