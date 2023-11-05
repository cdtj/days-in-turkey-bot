package country

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/i18n"
)

type Service interface {
	Info(ctx context.Context, l *i18n.Locale, c *model.Country) string
}
