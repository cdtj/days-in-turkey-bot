package http

import (
	"cdtj.io/days-in-turkey-bot/entity/user"
	httpserver "cdtj.io/days-in-turkey-bot/http-server"
)

// used for debug purposes
func RegisterHTTPEndpoints(router httpserver.HttpServerRouter, uc user.Usecase) {
	h := NewUserHttpHandler(uc)
	router.HandleFunc("/user/info", h.info)
	router.HandleFunc("/user/calc", h.calculateTrip)
	router.HandleFunc("/user/country", h.updateCountry)
	router.HandleFunc("/user/lang", h.updateLang)
}
