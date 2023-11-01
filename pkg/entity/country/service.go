package country

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
	"cdtj.io/days-in-turkey-bot/service/l10n"
)

type Service interface {
	Info(ctx context.Context, l *l10n.Locale, c *model.Country) string
}
