package http

import (
	"cdtj.io/days-in-turkey-bot/entity/country"
	"github.com/go-chi/chi/v5"
)

// RegisterHTTPEndpointsChi deprecated cuz echo is more handy in my scenario
func RegisterHTTPEndpointsChi(router *chi.Mux, uc country.Usecase) {
	h := NewCountryHttpHandlerChi(uc)
	router.HandleFunc("/country/get/{countryID}", h.getCountry)
}
