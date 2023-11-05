package http

import (
	"net/http"
	"strconv"

	"cdtj.io/days-in-turkey-bot/entity/user"
	"cdtj.io/days-in-turkey-bot/model"
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
	userID, err := strconv.ParseInt(r.URL.Query().Get("userID"), 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
	}
	user, err := h.usecase.Get(r.Context(), userID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
		return
	}
	resp, err := h.usecase.GetInfo(r.Context(), user)
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
	userID, err := strconv.ParseInt(r.URL.Query().Get("userID"), 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
	}
	input := new(CalculateTripInput)

	if err := render.DecodeJSON(r.Body, input); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
		return
	}

	user, err := h.usecase.Get(r.Context(), userID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
		return
	}
	resp, err := h.usecase.GetTrip(r.Context(), user, input.Dates)
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
	userID, err := strconv.ParseInt(r.URL.Query().Get("userID"), 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
	}
	input := new(UpdateLangInput)

	if err := render.DecodeJSON(r.Body, input); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
		return
	}
	user, err := h.usecase.Get(r.Context(), userID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
		return
	}
	err = h.usecase.UpdateLanguage(r.Context(), user, input.Lang)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, nil)
}

type UpdateCountryInput struct {
	Code          string `json:"code"`
	Flag          string `json:"flag"`
	DaysContinual int    `json:"daysContinual"`
	DaysLimit     int    `json:"daysLimit"`
	ResetInterval int    `json:"resetInterval"`
}

func (h *UserHttpHandler) updateCountry(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.URL.Query().Get("userID"), 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
	}
	input := new(UpdateCountryInput)

	if err := render.DecodeJSON(r.Body, input); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
		return
	}
	user, err := h.usecase.Get(r.Context(), userID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorUserResponse{err.Error()})
		return
	}
	// this is test handler for internal use so we don't make any data validation
	// use at your own risk
	err = h.usecase.UpdateCountry(r.Context(), user, model.NewCountry(input.Code, input.Flag, input.DaysContinual, input.DaysLimit, input.ResetInterval))
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
