package http

import (
	"cdtj.io/days-in-turkey-bot/entity/user"
	"github.com/gorilla/mux"
)

func RegisterHTTPEndpoints(router *mux.Router, uc user.Usecase) {
	h := NewUserHandler(uc)
	router.HandleFunc("/user/get", h.getUser)
}
