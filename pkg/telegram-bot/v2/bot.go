package telegrambot

import (
	"context"
	"errors"
	"log/slog"

	"cdtj.io/days-in-turkey-bot/model"
	tgapi "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type TelegramBot struct {
	bot     *tgapi.Bot
	options *TelegramBotOptions
}

func NewTelegramBot(token string, options *TelegramBotOptions) *TelegramBot {
	b, err := tgapi.New(token, options.apiOptions()...)
	if err != nil {
		slog.Error("telegram-bot", "msg", "unable to create new bot", "err", err)
		return nil
	}
	return &TelegramBot{
		bot:     b,
		options: options,
	}
}

var (
	ErrBotNotReady = errors.New("not initialized or already stopped")
)

func (t *TelegramBot) Serve(ctx context.Context) error {
	slog.Info("telegram-bot", "status", "starting")
	if t == nil || t.bot == nil {
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

func (t *TelegramBot) SetCommands(ctx context.Context, commands []*model.TelegramBotCommand, languageCode string) error {
	tgcmds := make([]models.BotCommand, 0, len(commands))
	for _, command := range commands {
		tgcmds = append(tgcmds, models.BotCommand{
			Command:     command.Command,
			Description: command.Caption,
		})
	}
	_, err := t.bot.SetMyCommands(ctx, &tgapi.SetMyCommandsParams{
		Commands:     tgcmds,
		LanguageCode: languageCode,
	})
	return err
}

func (t *TelegramBot) SetDescription(ctx context.Context, description, about, languageCode string) error {
	if _, err := t.bot.SetMyDescription(ctx, &tgapi.SetMyDescriptionParams{
		Description:  description,
		LanguageCode: languageCode,
	}); err != nil {
		return err
	}
	if _, err := t.bot.SetMyShortDescription(ctx, &tgapi.SetMyShortDescriptionParams{
		ShortDescription: about,
		LanguageCode:     languageCode,
	}); err != nil {
		return err
	}
	return nil
}
