package bot

import (
	"context"
)

type Service interface {
	Reply(ctx context.Context, chatID int64, text string) error
}

type ServiceSender interface {
	Send(ctx context.Context, chatID int64, text string) error
}
