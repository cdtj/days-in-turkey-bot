package telegrambot

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"cdtj.io/days-in-turkey-bot/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		debug:   false,
	}
}

// SetDebug sets bot into the debug mode, must be called before Serve()
func (t *TelegramBot) SetDebug(debug bool) {
	t.debug = debug
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

	wh, err := tgbotapi.NewWebhook(t.webhook)
	if err != nil {
		return err
	}
	resp, err := bot.Request(wh)
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

	t.bot.StopReceivingUpdates()
	t.bot = nil
	slog.Info("bot stopped")
}

var (
	ErrBotNotReady = errors.New("not not started or already stopped")
)

func (t *TelegramBot) Send(ctx context.Context, chatID int64, text string, commands []*model.TelegramBotCommandRow) error {
	if t.bot != nil {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = tgbotapi.ModeMarkdown
		if commands != nil {
			msg.ReplyMarkup = replyMarkup(ctx, commands)
		}
		if _, err := t.bot.Send(msg); err != nil {
			return err
		}
	}
	return ErrBotNotReady
}

func replyMarkup(ctx context.Context, rows []*model.TelegramBotCommandRow) tgbotapi.InlineKeyboardMarkup {
	ikbs := make([]tgbotapi.InlineKeyboardButton, 0)
	for _, row := range rows {
		ikrs := make([]tgbotapi.InlineKeyboardButton, 0)
		for _, command := range row.Commands {
			ikrs = append(ikrs, tgbotapi.NewInlineKeyboardButtonData(command.Caption, command.Command))
		}
		ikbs = append(ikbs, ikrs...)
	}
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(ikbs)
	return inlineKeyboard
}
