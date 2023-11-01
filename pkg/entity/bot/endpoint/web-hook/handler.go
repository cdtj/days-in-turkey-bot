package webhook

import (
	"log/slog"
	"net/http"
	"strconv"

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
	if msg == nil {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, nil)
		return
	}

	slog.Debug("incoming",
		"UpdateID", update.UpdateID,
		slog.Group("Message",
			"From", msg.From,
			"Text", msg.Text,
			"Entities", msg.Entities,
		),
		slog.Group("Command",
			"IsCommand", msg.IsCommand(),
			"Command", msg.Command(),
			"CommandArguments", msg.CommandArguments(),
			"CommandWithAt", msg.CommandWithAt(),
		),
	)

	var resp string
	userID := strconv.Itoa(msg.From.ID)
	if msg.IsCommand() {
		switch msg.Command() {
		case BotWebhookCountry:
			resp, err = h.usecase.UpdateCountry(r.Context(), userID, msg.CommandArguments())
		case BotWebhookLanguage:
			resp, err = h.usecase.UpdateLang(r.Context(), userID, msg.CommandArguments())
		case BotWebhookStart:
			resp = "hi"
		}
	} else {
		resp, err = h.usecase.CalculateTrip(r.Context(), userID, msg.Text)
	}
	if err != nil {
		h.usecase.Reply(r.Context(), msg.Chat.ID, err.Error())
		return
	}
	h.usecase.Reply(r.Context(), msg.Chat.ID, resp)
	slog.Info("result", "resp", resp, "err", err)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, nil)
}

type ErrorBotResponse struct {
	Error string `json:"error"`
}

type BotResponse struct {
	Response string `json:"response"`
}
