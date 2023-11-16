package telegrambot

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
)

type TelegramBotv2 struct {
	bot *bot.Bot
}

func NewTelegramBotv2(token string) *TelegramBotv2 {
	b, err := bot.New(token)
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
	ok, err := t.bot.DeleteWebhook(ctx, &bot.DeleteWebhookParams{DropPendingUpdates: true})
	slog.Info("telegram-bot", "msg", "cleanup old hooks", "ok", ok, "err", err)
	t.bot.Start(ctx)
	return nil
}

func (t *TelegramBotv2) Shutdown(ctx context.Context) {
	// there is no designated shutdown
	t.bot = nil
	slog.Info("telegram-bot", "status", "stopped")
}

func (t *TelegramBotv2) RegisterHandlerExactMessage(command string, handler bot.HandlerFunc) {
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, command, bot.MatchTypeExact, handler)
}

func (t *TelegramBotv2) RegisterHandlerPrefixMessage(command string, handler bot.HandlerFunc) {
	t.bot.RegisterHandler(bot.HandlerTypeMessageText, command, bot.MatchTypePrefix, handler)
}

func (t *TelegramBotv2) RegisterHandlerExactCb(command string, handler bot.HandlerFunc) {
	t.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, command, bot.MatchTypeExact, handler)
}

func (t *TelegramBotv2) RegisterHandlerPrefixCb(command string, handler bot.HandlerFunc) {
	t.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, command, bot.MatchTypePrefix, handler)
}
