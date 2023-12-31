package http

import (
	"net/http"
	"strconv"

	"cdtj.io/days-in-turkey-bot/entity/user"
	"cdtj.io/days-in-turkey-bot/model"
	"github.com/labstack/echo/v4"
)

type UserHttpHandlerEcho struct {
	usecase user.Usecase
}

func NewUserHttpHandlerEcho(usecase user.Usecase) *UserHttpHandlerEcho {
	return &UserHttpHandlerEcho{
		usecase: usecase,
	}
}

func (h *UserHttpHandlerEcho) info(c echo.Context) error {
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		return err
	}
	user, err := h.usecase.Get(c.Request().Context(), userID)
	if err != nil {
		return err
	}
	resp := h.usecase.GetInfo(c.Request().Context(), user)
	return c.JSON(http.StatusOK, &UserResponse{resp})
}

type CalculateTripInput struct {
	Dates string `json:"dates"`
}

func (h *UserHttpHandlerEcho) calculateTrip(c echo.Context) error {
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		return err
	}
	input := new(CalculateTripInput)
	if err := c.Bind(input); err != nil {
		return err
	}
	user, err := h.usecase.Get(c.Request().Context(), userID)
	if err != nil {
		return err
	}
	resp, err := h.usecase.CalculateTrip(c.Request().Context(), user, input.Dates)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &UserResponse{resp})
}

type UpdateLangInput struct {
	Lang string `json:"lang"`
}

func (h *UserHttpHandlerEcho) updateLang(c echo.Context) error {
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		return err
	}
	input := new(UpdateLangInput)
	if err := c.Bind(input); err != nil {
		return err
	}
	user, err := h.usecase.Get(c.Request().Context(), userID)
	if err != nil {
		return err
	}
	if err := h.usecase.UpdateLanguage(c.Request().Context(), user, input.Lang); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

type UpdateCountryInput struct {
	Name          string `json:"name"`
	Code          string `json:"code"`
	Flag          string `json:"flag"`
	DaysContinual int    `json:"daysContinual"`
	DaysLimit     int    `json:"daysLimit"`
	ResetInterval int    `json:"resetInterval"`
	VisaFree      bool   `json:"visaFree"`
}

func (h *UserHttpHandlerEcho) updateCountry(c echo.Context) error {
	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		return err
	}
	input := new(UpdateCountryInput)
	if err := c.Bind(input); err != nil {
		return err
	}
	user, err := h.usecase.Get(c.Request().Context(), userID)
	if err != nil {
		return err
	}
	// this is test handler for internal use so we don't make any data validation
	// use at your own risk
	if err := h.usecase.UpdateCountry(c.Request().Context(), user, model.NewCountry(input.Code, input.Flag, input.Name, input.DaysContinual, input.DaysLimit, input.ResetInterval, input.VisaFree)); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}

type ErrorUserResponse struct {
	Error string `json:"error"`
}

type UserResponse struct {
	Response string `json:"response"`
}
