package webhook

import (
	"log/slog"
	"net/http"
	"strings"

	"cdtj.io/days-in-turkey-bot/entity/bot"
	"github.com/go-chi/render"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotWebhookHandlerChi struct {
	usecase bot.Usecase
}

func NewBotWebhookHandlerChi(usecase bot.Usecase) *BotWebhookHandlerChi {
	return &BotWebhookHandlerChi{
		usecase: usecase,
	}
}

func (h *BotWebhookHandlerChi) webhook(w http.ResponseWriter, r *http.Request) {
	update := new(tgbotapi.Update)
	err := render.DecodeJSON(r.Body, update)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &ErrorBotResponse{err.Error()})
		return
	}
	msg := update.Message
	cb := update.CallbackQuery
	if msg == nil && cb == nil {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, nil)
		return
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
				err = h.usecase.Prompt(r.Context(), chatID, userID, msg.Command())
			case BotWebhookStart:
				err = h.usecase.Welcome(r.Context(), chatID, userID, msg.From.LanguageCode)
			}
		} else {
			err = h.usecase.CalculateTrip(r.Context(), chatID, userID, msg.Text)
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
				err = h.usecase.UpdateCountry(r.Context(), chatID, userID, inputArr[1], -1, -1, -1)
			case BotWebhookLanguage:
				err = h.usecase.UpdateLanguage(r.Context(), chatID, userID, inputArr[1])
			}
		}
	}

	if err != nil {
		slog.Info("result", "err", err)
		if chatID > 0 {
			h.usecase.Send(r.Context(), chatID, err.Error(), nil)
		}
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, nil)
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
