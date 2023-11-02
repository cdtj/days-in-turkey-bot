package bot

import (
	"context"

	"cdtj.io/days-in-turkey-bot/model"
)

type Service interface {
	Send(ctx context.Context, chatID int64, text string, replyMarkup []*model.TelegramBotCommandRow) error
}
