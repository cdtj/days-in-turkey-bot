package http

import (
	"net/http"

	"cdtj.io/days-in-turkey-bot/entity/country"
	"cdtj.io/days-in-turkey-bot/service/i18n"
	"github.com/labstack/echo/v4"
)

type CountryHttpHandlerEcho struct {
	usecase country.Usecase
}

func NewCountryHttpHandlerEcho(usecase country.Usecase) *CountryHttpHandlerEcho {
	return &CountryHttpHandlerEcho{
		usecase: usecase,
	}
}

func (h *CountryHttpHandlerEcho) getCountry(c echo.Context) error {
	countryID := c.Param("countryID")
	input := new(GetCountryInput)

	if err := c.Bind(input); err != nil {
		return err
	}

	country, err := h.usecase.Get(c.Request().Context(), countryID)
	if err != nil {
		return err
	}

	language, err := i18n.LanguageLookup(input.Lang)
	if err != nil {
		return err
	}

	resp, err := h.usecase.GetInfo(c.Request().Context(), language, country)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &CountryResponse{resp})
}
