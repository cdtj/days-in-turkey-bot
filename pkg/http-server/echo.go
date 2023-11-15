package httpserver

import (
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewEchoRouter() *echo.Echo {
	e := echo.New()

	// e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	// logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				slog.Info(c.RealIP(),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				slog.Error(c.RealIP(),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	return e
}
