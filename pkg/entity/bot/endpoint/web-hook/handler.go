package webhook

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"cdtj.io/days-in-turkey-bot/entity/bot"
	"github.com/go-chi/render"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type BotWebhookHandler struct {
	usecase bot.Usecase
}

func NewBotWebhookHandler(usecase bot.Usecase) *BotWebhookHandler {
	return &BotWebhookHandler{
		usecase: usecase,
	}
}

const (
	BotWebhookCountry  = "country"
	BotWebhookLanguage = "language"
	BotWebhookStart    = "start"
)

func (h *BotWebhookHandler) webhook(w http.ResponseWriter, r *http.Request) {
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
		userID := strconv.Itoa(msg.From.ID)
		if msg.IsCommand() {
			switch msg.Command() {
			case BotWebhookCountry:
				err = h.usecase.UpdateCountry(r.Context(), msg.Chat.ID, userID, msg.CommandArguments())
			case BotWebhookLanguage:
				err = h.usecase.UpdateLang(r.Context(), msg.Chat.ID, userID, msg.CommandArguments())
			case BotWebhookStart:
				err = h.usecase.Welcome(r.Context(), msg.Chat.ID, userID, msg.From.LanguageCode)
			}
		} else {
			err = h.usecase.CalculateTrip(r.Context(), msg.Chat.ID, userID, msg.Text)
		}
	case cb != nil:
		slog.Info("incoming callback",
			"UpdateID", update.UpdateID,
			slog.Group("Message",
				"From", cb.From,
				"Data", cb.Data,
			),
		)
		userID := strconv.Itoa(cb.From.ID)
		inputArr := strings.Split(cb.Data, " ")
		if len(inputArr) == 2 {
			switch inputArr[0] {
			case BotWebhookCountry:
				err = h.usecase.UpdateCountry(r.Context(), cb.Message.Chat.ID, userID, inputArr[1])
			case BotWebhookLanguage:
				err = h.usecase.UpdateLang(r.Context(), cb.Message.Chat.ID, userID, inputArr[1])
			}
		}
	}

	if err != nil {
		slog.Info("result", "err", err)
		h.usecase.Send(r.Context(), msg.Chat.ID, err.Error(), nil)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, nil)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, nil)
}

type ErrorBotResponse struct {
	Error string `json:"error"`
}

type BotResponse struct {
	Response string `json:"response"`
}
