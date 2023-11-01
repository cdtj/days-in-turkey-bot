package http

import (
	"net/http"

	"cdtj.io/days-in-turkey-bot/entity/country"
	"github.com/go-chi/render"
)

type CountryHttpHandler struct {
	usecase country.Usecase
}

func NewCountryHttpHandler(usecase country.Usecase) *CountryHttpHandler {
	return &CountryHttpHandler{
		usecase: usecase,
	}
}

func (h *CountryHttpHandler) getCountry(w http.ResponseWriter, r *http.Request) {
	countryID := r.URL.Query().Get("countryID")
	resp, err := h.usecase.Info(r.Context(), countryID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorCountryResponse{err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &CountryResponse{resp})
}

type ErrorCountryResponse struct {
	Error string `json:"error"`
}

type CountryResponse struct {
	Response string `json:"response"`
}
