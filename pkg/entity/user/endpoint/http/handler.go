package http

import (
	"net/http"

	"cdtj.io/days-in-turkey-bot/entity/user"
	"github.com/go-chi/render"
)

type UserHttpHandler struct {
	usecase user.Usecase
}

func NewUserHttpHandler(usecase user.Usecase) *UserHttpHandler {
	return &UserHttpHandler{
		usecase: usecase,
	}
}

func (h *UserHttpHandler) info(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	resp, err := h.usecase.Info(r.Context(), userID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &UserResponse{resp})
}

type CalculateTripInput struct {
	Dates string `json:"dates"`
}

func (h *UserHttpHandler) calculateTrip(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	input := new(CalculateTripInput)

	if err := render.DecodeJSON(r.Body, input); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
		return
	}

	resp, err := h.usecase.CalculateTrip(r.Context(), userID, input.Dates)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &UserResponse{resp})
}

type UpdateLangInput struct {
	Lang string `json:"lang"`
}

func (h *UserHttpHandler) updateLang(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	input := new(UpdateLangInput)

	if err := render.DecodeJSON(r.Body, input); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
		return
	}
	err := h.usecase.UpdateLang(r.Context(), userID, input.Lang)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, nil)
}

type UpdateCountryInput struct {
	Country string `json:"country"`
}

func (h *UserHttpHandler) updateCountry(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	input := new(UpdateCountryInput)

	if err := render.DecodeJSON(r.Body, input); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
		return
	}
	err := h.usecase.UpdateCountry(r.Context(), userID, input.Country)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, nil)
}

type ErrorUserResponse struct {
	Error string `json:"error"`
}

type UserResponse struct {
	Response string `json:"response"`
}
