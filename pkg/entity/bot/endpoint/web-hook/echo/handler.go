package webhook

import (
	"log/slog"
	"net/http"
	"strings"

	"cdtj.io/days-in-turkey-bot/entity/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/labstack/echo/v4"
)

type BotWebhookHandlerEcho struct {
	usecase bot.Usecase
}

func NewBotWebhookHandlerEcho(usecase bot.Usecase) *BotWebhookHandlerEcho {
	return &BotWebhookHandlerEcho{
		usecase: usecase,
	}
}

func (h *BotWebhookHandlerEcho) webhook(c echo.Context) error {
	update := new(tgbotapi.Update)
	err := c.Bind(update)
	if err != nil {
		return err
	}
	msg := update.Message
	cb := update.CallbackQuery
	if msg == nil && cb == nil {
		return c.JSON(http.StatusOK, nil)
	}
	var chatID, userID int64

	switch {
	case msg != nil:
		slog.Info("incoming message",
			"UpdateID", update.UpdateID,
			slog.Group("Message",
				"From", msg.From,
				"Text", msg.Text,
				"Entities", msg.Entities,
			),
			slog.Group("User",
				"Name", msg.From.String(),
				"Lang", msg.From.LanguageCode,
			),
			slog.Group("Command",
				"IsCommand", msg.IsCommand(),
				"Command", msg.Command(),
				"CommandArguments", msg.CommandArguments(),
				"CommandWithAt", msg.CommandWithAt(),
			),
		)
		chatID = msg.Chat.ID
		userID = msg.From.ID
		if msg.IsCommand() {
			switch msg.Command() {
			case BotWebhookCountry, BotWebhookLanguage, BotWebhookContribute, BotWebhookTrip:
				err = h.usecase.Prompt(c.Request().Context(), chatID, userID, msg.Command())
			case BotWebhookStart:
				err = h.usecase.Welcome(c.Request().Context(), chatID, userID, msg.From.LanguageCode)
			}
		} else {
			err = h.usecase.CalculateTrip(c.Request().Context(), chatID, userID, msg.Text)
		}
	case cb != nil:
		slog.Info("incoming callback",
			"UpdateID", update.UpdateID,
			slog.Group("Message",
				"From", cb.From,
				"Data", cb.Data,
			),
		)
		chatID = cb.Message.Chat.ID
		userID = cb.From.ID
		inputArr := strings.Split(cb.Data, " ")
		if len(inputArr) == 2 {
			switch inputArr[0] {
			case BotWebhookCountry:
				err = h.usecase.UpdateCountry(c.Request().Context(), chatID, userID, inputArr[1], -1, -1, -1)
			case BotWebhookLanguage:
				err = h.usecase.UpdateLanguage(c.Request().Context(), chatID, userID, inputArr[1])
			}
		}
	}

	if err != nil {
		slog.Info("result", "err", err)
		if chatID > 0 {
			h.usecase.Send(c.Request().Context(), chatID, err.Error(), nil)
		}
	}
	return c.JSON(http.StatusOK, nil)
}

const (
	BotWebhookCountry    = "country"
	BotWebhookLanguage   = "language"
	BotWebhookContribute = "contribute"
	BotWebhookTrip       = "trip"
	BotWebhookStart      = "start"
)

type ErrorBotResponse struct {
	Error string `json:"error"`
}

type BotResponse struct {
	Response string `json:"response"`
}
