package http

import (
	"cdtj.io/days-in-turkey-bot/entity/user"
	"github.com/go-chi/chi/v5"
)

// RegisterHTTPEndpointsChi deprecated cuz echo is more handy in my scenario
func RegisterHTTPEndpointsChi(router *chi.Mux, uc user.Usecase) {
	h := NewUserHttpHandlerChi(uc)
	router.HandleFunc("/user/info/{userID}", h.info)
	router.HandleFunc("/user/calc/{userID}", h.calculateTrip)
	router.HandleFunc("/user/country/{userID}", h.updateCountry)
	router.HandleFunc("/user/lang/{userID}", h.updateLang)
}
