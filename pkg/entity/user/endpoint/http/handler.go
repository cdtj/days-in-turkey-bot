package http

import (
	"net/http"

	"cdtj.io/days-in-turkey-bot/entity/user"
)

type UserHandler struct {
	usecase user.Usecase
}

func NewUserHandler(usecase user.Usecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}

func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	resp, err := h.usecase.Get(r.Context(), "id")
	if err != nil {
		return
	}
	w.Write([]byte(resp))
}
