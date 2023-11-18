package telegrambot

import (
	"context"
	"errors"
	"log/slog"

	tgapi "github.com/go-telegram/bot"
)

type TelegramBot struct {
	bot *tgapi.Bot
}

func NewTelegramBot(token string, options []tgapi.Option) *TelegramBot {
	b, err := tgapi.New(token, options...)
	if err != nil {
		return nil
	}
	return &TelegramBot{
		bot: b,
	}
}

var (
	ErrBotNotReady = errors.New("not not started or already stopped")
)

func (t *TelegramBot) Serve(ctx context.Context) error {
	slog.Info("telegram-bot", "status", "starting")
	if t.bot == nil {
		return ErrBotNotReady
	}
	ok, err := t.bot.DeleteWebhook(ctx, &tgapi.DeleteWebhookParams{DropPendingUpdates: true})
	slog.Info("telegram-bot", "msg", "cleanup old hooks", "ok", ok, "err", err)
	t.bot.Start(ctx)
	return nil
}

func (t *TelegramBot) Shutdown(ctx context.Context) {
	// there is no designated shutdown

	// you can restart bot without re-init, so do not delete it on stop
	// t.bot = nil
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
