package telegrambot

import (
	"context"
	"log/slog"

	tgapi "github.com/go-telegram/bot"
)

type TelegramBotv2 struct {
	bot *tgapi.Bot
}

func NewTelegramBotv2(token string, options []tgapi.Option) *TelegramBotv2 {
	b, err := tgapi.New(token, options...)
	if err != nil {
		return nil
	}
	return &TelegramBotv2{
		bot: b,
	}
}

func (t *TelegramBotv2) Serve(ctx context.Context) error {
	slog.Info("telegram-bot", "status", "starting")
	if t.bot == nil {
		return ErrBotNotReady
	}
	ok, err := t.bot.DeleteWebhook(ctx, &tgapi.DeleteWebhookParams{DropPendingUpdates: true})
	slog.Info("telegram-bot", "msg", "cleanup old hooks", "ok", ok, "err", err)
	t.bot.Start(ctx)
	return nil
}

func (t *TelegramBotv2) Shutdown(ctx context.Context) {
	// there is no designated shutdown
	t.bot = nil
	slog.Info("telegram-bot", "status", "stopped")
}

func BindHandlerExactMessage(command string, handler tgapi.HandlerFunc) tgapi.Option {
	return tgapi.WithMessageTextHandler(command, tgapi.MatchTypeExact, handler)
}

func BindHandlerPrefixMessage(command string, handler tgapi.HandlerFunc) tgapi.Option {
	return tgapi.WithMessageTextHandler(command, tgapi.MatchTypePrefix, handler)
}

func BindHandlerExactCb(command string, handler tgapi.HandlerFunc) tgapi.Option {
	return tgapi.WithCallbackQueryDataHandler(command, tgapi.MatchTypeExact, handler)
}

func BindHandlerPrefixCb(command string, handler tgapi.HandlerFunc) tgapi.Option {
	return tgapi.WithCallbackQueryDataHandler(command, tgapi.MatchTypePrefix, handler)
}

func BindHandlerDefault(handler tgapi.HandlerFunc) tgapi.Option {
	return tgapi.WithDefaultHandler(handler)
}
