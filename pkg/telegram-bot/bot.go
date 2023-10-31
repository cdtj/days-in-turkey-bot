package telegrambot

import (
	"context"
	"errors"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramBot struct {
	token   string
	webhook string
	debug   bool
	bot     *tgbotapi.BotAPI
	sender  func(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

func NewTelegramBot(token, webhook string) *TelegramBot {
	return &TelegramBot{
		token:   token,
		webhook: webhook,
		debug:   true,
	}
}

func (t *TelegramBot) Serve(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(t.token)
	if err != nil {
		return err
	}
	t.bot = bot
	t.bot.Debug = t.debug
	t.sender = t.bot.Send
	slog.Info("authorized", "account", t.bot.Self.UserName)

	resp, err := bot.SetWebhook(tgbotapi.NewWebhook(t.webhook))
	if err != nil {
		return err
	}
	slog.Info("webhook deployed", "success", resp.Result)

	<-ctx.Done()
	return nil
}

func (t *TelegramBot) Shutdown(ctx context.Context) {
	shutdownCtx, shutdownStopCtx := context.WithTimeout(ctx, 30*time.Second)
	defer shutdownStopCtx()

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			slog.Error("unable to gracefully stop telegram bot", "error", shutdownCtx.Err())
			return
		}
	}()

	resp, err := t.bot.RemoveWebhook()
	if err != nil {
		slog.Error("unable to remove webhook", "error", err)
		return
	}
	t.bot = nil
	slog.Info("webhook removed", "result", resp.Result)
}

var (
	ErrBotNotReady = errors.New("not not started or already stopped")
)

func (t *TelegramBot) Send(ctx context.Context, chatID int64, text string) error {
	if t.bot != nil {
		msg := tgbotapi.NewMessage(chatID, text)
		if _, err := t.bot.Send(msg); err != nil {
			return err
		}
	}
	return ErrBotNotReady
}
