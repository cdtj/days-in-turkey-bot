package http

import (
	"net/http"

	"cdtj.io/days-in-turkey-bot/entity/country"
	"cdtj.io/days-in-turkey-bot/service/i18n"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type CountryHttpHandlerChi struct {
	usecase country.Usecase
}

func NewCountryHttpHandlerChi(usecase country.Usecase) *CountryHttpHandlerChi {
	return &CountryHttpHandlerChi{
		usecase: usecase,
	}
}

func (h *CountryHttpHandlerChi) getCountry(w http.ResponseWriter, r *http.Request) {
	countryID := chi.URLParam(r, "countryID")
	input := new(GetCountryInput)

	if err := render.DecodeJSON(r.Body, input); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorCountryResponse{err.Error()})
		return
	}

	country, err := h.usecase.Get(r.Context(), countryID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorCountryResponse{err.Error()})
		return
	}

	language, err := i18n.LanguageLookup(input.Lang)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorCountryResponse{err.Error()})
		return
	}

	resp, err := h.usecase.GetInfo(r.Context(), language, country)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorCountryResponse{err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &CountryResponse{resp})
}
