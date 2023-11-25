package tghandler

import (
	"context"
	"log/slog"

	"cdtj.io/days-in-turkey-bot/entity/bot/v2"
	tgapi "github.com/go-telegram/bot"
	tgmodel "github.com/go-telegram/bot/models"
)

// This place contains the only public-exposable Handlers
// all other Handlers must be used for debug-only purposes

type BotHandler struct {
	usecase bot.Usecase
}

func NewBotHandler(usecase bot.Usecase) *BotHandler {
	return &BotHandler{
		usecase: usecase,
	}
}

// welcome handles user greetings
//   - input: [update.Message]
//   - usecase: [usecase.Welcome], [usecase.Me]
func (h *BotHandler) welcome(ctx context.Context, b *tgapi.Bot, update *tgmodel.Update) {
	tgmsg := update.Message
	if tgmsg == nil || tgmsg.From == nil {
		slog.Error("Unexpected message with empty message", "method", "welcome", "update", update)
		return
	}
	chatID := tgmsg.Chat.ID
	msg := h.usecase.Welcome(ctx, chatID, tgmsg.From.LanguageCode)
	if _, err := b.SendMessage(ctx, &tgapi.SendMessageParams{
		ChatID:      chatID,
		Text:        msg.Text,
		ParseMode:   tgmodel.ParseModeMarkdown,
		ReplyMarkup: msg.Markup,
	}); err != nil {
		slog.Error("SendMessage failed", "method", "welcome", "chatID", chatID, "text", msg.Text, "err", err)
	}
	meMsg := h.usecase.Me(ctx, chatID)
	if _, err := b.SendMessage(ctx, &tgapi.SendMessageParams{
		ChatID:      chatID,
		Text:        meMsg.Text,
		ParseMode:   tgmodel.ParseModeMarkdown,
		ReplyMarkup: meMsg.Markup,
	}); err != nil {
		slog.Error("SendMessage failed", "method", "me", "chatID", chatID, "text", meMsg.Text, "err", err)
	}
}

// country handles user country-related info
//   - input: [update.Message]
//   - usecase: [usecase.Country], [usecase.Hint]
func (h *BotHandler) country(ctx context.Context, b *tgapi.Bot, update *tgmodel.Update) {
	tgmsg := update.Message
	if tgmsg == nil {
		slog.Error("Unexpected message with empty message", "method", "country", "update", update)
		return
	}
	chatID := tgmsg.Chat.ID
	msg := h.usecase.Country(ctx, chatID)
	if _, err := b.SendMessage(ctx, &tgapi.SendMessageParams{
		ChatID:      chatID,
		Text:        msg.Text,
		ParseMode:   tgmodel.ParseModeMarkdown,
		ReplyMarkup: msg.Markup,
	}); err != nil {
		slog.Error("SendMessage failed", "method", "country", "chatID", chatID, "text", msg.Text, "err", err)
	}

	hint := h.usecase.Hint(ctx, chatID, "UserCountryCustom")
	if _, err := b.SendMessage(ctx, &tgapi.SendMessageParams{
		ChatID:      chatID,
		Text:        hint.Text,
		ParseMode:   tgmodel.ParseModeMarkdown,
		ReplyMarkup: hint.Markup,
	}); err != nil {
		slog.Error("SendMessage failed", "method", "country", "chatID", chatID, "text", hint.Text, "err", err)
	}
}

// language handles user language-related info
//   - input: [update.Message]
//   - usecase: [usecase.Language]
func (h *BotHandler) language(ctx context.Context, b *tgapi.Bot, update *tgmodel.Update) {
	tgmsg := update.Message
	if tgmsg == nil {
		slog.Error("Unexpected message with empty message", "method", "language", "update", update)
		return
	}
	chatID := tgmsg.Chat.ID
	msg := h.usecase.Language(ctx, chatID)
	if _, err := b.SendMessage(ctx, &tgapi.SendMessageParams{
		ChatID:      chatID,
		Text:        msg.Text,
		ParseMode:   tgmodel.ParseModeMarkdown,
		ReplyMarkup: msg.Markup,
	}); err != nil {
		slog.Error("SendMessage failed", "method", "language", "chatID", chatID, "text", msg.Text, "err", err)
	}
}

// contribute handles user contribution-related info
//   - input: [update.Message]
//   - usecase: [usecase.Contribute]
func (h *BotHandler) contribute(ctx context.Context, b *tgapi.Bot, update *tgmodel.Update) {
	tgmsg := update.Message
	if tgmsg == nil {
		slog.Error("Unexpected message with empty message", "method", "contribute", "update", update)
		return
	}
	chatID := tgmsg.Chat.ID
	msg := h.usecase.Contribute(ctx, chatID)
	if _, err := b.SendMessage(ctx, &tgapi.SendMessageParams{
		ChatID:      chatID,
		Text:        msg.Text,
		ParseMode:   tgmodel.ParseModeMarkdown,
		ReplyMarkup: msg.Markup,
	}); err != nil {
		slog.Error("SendMessage failed", "method", "contribute", "chatID", chatID, "text", msg.Text, "err", err)
	}
}

// trip handles user trip-related info
//   - input: [update.Message]
//   - usecase: [usecase.Trip]
func (h *BotHandler) trip(ctx context.Context, b *tgapi.Bot, update *tgmodel.Update) {
	tgmsg := update.Message
	if tgmsg == nil {
		slog.Error("Unexpected message with empty message", "method", "trip", "update", update)
		return
	}
	chatID := tgmsg.Chat.ID
	msg := h.usecase.Trip(ctx, chatID)
	if _, err := b.SendMessage(ctx, &tgapi.SendMessageParams{
		ChatID:      chatID,
		Text:        msg.Text,
		ParseMode:   tgmodel.ParseModeMarkdown,
		ReplyMarkup: msg.Markup,
	}); err != nil {
		slog.Error("SendMessage failed", "method", "trip", "chatID", chatID, "text", msg.Text, "err", err)
	}
}

// me handles user user-related info
//   - input: [update.Message]
//   - usecase: [usecase.Me]
func (h *BotHandler) me(ctx context.Context, b *tgapi.Bot, update *tgmodel.Update) {
	tgmsg := update.Message
	if tgmsg == nil {
		slog.Error("Unexpected message with empty message", "method", "me", "update", update)
		return
	}
	chatID := tgmsg.Chat.ID
	msg := h.usecase.Me(ctx, chatID)
	if _, err := b.SendMessage(ctx, &tgapi.SendMessageParams{
		ChatID:      chatID,
		Text:        msg.Text,
		ParseMode:   tgmodel.ParseModeMarkdown,
		ReplyMarkup: msg.Markup,
	}); err != nil {
		slog.Error("SendMessage failed", "method", "me", "chatID", chatID, "text", msg.Text, "err", err)
	}
}

// feedback handles user feedback-related info
//   - input: [update.Message]
//   - usecase: [usecase.Feedback]
func (h *BotHandler) feedback(ctx context.Context, b *tgapi.Bot, update *tgmodel.Update) {
	tgmsg := update.Message
	if tgmsg == nil {
		slog.Error("Unexpected message with empty message", "method", "me", "update", update)
		return
	}
	chatID := tgmsg.Chat.ID
	msg := h.usecase.Feedback(ctx, chatID)
	if _, err := b.SendMessage(ctx, &tgapi.SendMessageParams{
		ChatID:      chatID,
		Text:        msg.Text,
		ParseMode:   tgmodel.ParseModeMarkdown,
		ReplyMarkup: msg.Markup,
	}); err != nil {
		slog.Error("SendMessage failed", "method", "me", "chatID", chatID, "text", msg.Text, "err", err)
	}
}

// updateCountry handles country changing request
//   - input: [update.Message], [update.CallbackQuery]
//   - usecase: [usecase.UpdateCountry]
func (h *BotHandler) updateCountry(ctx context.Context, b *tgapi.Bot, update *tgmodel.Update) {
	cb := update.CallbackQuery
	tgmsg := update.Message
	if cb == nil && tgmsg == nil {
		slog.Error("Unexpected callback with empty query", "method", "updateCountry", "update", update)
		return
	}
	var chatID int64
	var countryInput string
	switch {
	case cb != nil:
		chatID = cb.Sender.ID
		countryInput = cb.Data
	case tgmsg != nil:
		chatID = tgmsg.Chat.ID
		countryInput = tgmsg.Text
	}
	msg := h.usecase.UpdateCountry(ctx, chatID, countryInput)
	if _, err := b.SendMessage(ctx, &tgapi.SendMessageParams{
		ChatID:      chatID,
		Text:        msg.Text,
		ParseMode:   tgmodel.ParseModeMarkdown,
		ReplyMarkup: msg.Markup,
	}); err != nil {
		slog.Error("SendMessage failed", "method", "updateCountry", "chatID", chatID, "input", countryInput, "text", msg.Text, "err", err)
	}
}

// updateLanguage handles language changing request
//   - input: [update.CallbackQuery]
//   - usecase: [usecase.UpdateLanguage]
func (h *BotHandler) updateLanguage(ctx context.Context, b *tgapi.Bot, update *tgmodel.Update) {
	cb := update.CallbackQuery
	if cb == nil {
		slog.Error("Unexpected callback with empty query", "method", "updateLanguage", "update", update)
		return
	}
	chatID := cb.Sender.ID
	languageInput := cb.Data
	msg := h.usecase.UpdateLanguage(ctx, chatID, languageInput)
	if _, err := b.SendMessage(ctx, &tgapi.SendMessageParams{
		ChatID:      chatID,
		Text:        msg.Text,
		ParseMode:   tgmodel.ParseModeMarkdown,
		ReplyMarkup: msg.Markup,
	}); err != nil {
		slog.Error("SendMessage failed", "method", "updateLanguage", "chatID", chatID, "input", languageInput, "text", msg.Text, "err", err, "msg", msg.Text)
	}
}

// defaultMessage handles default request meant to be a trip calculation input
//   - input: [update.Message]
//   - usecase: [usecase.CalculateTrip]
func (h *BotHandler) defaultMessage(ctx context.Context, b *tgapi.Bot, update *tgmodel.Update) {
	tgmsg := update.Message
	if tgmsg == nil {
		slog.Error("Unexpected message with empty message", "method", "defaultMessage", "update", update)
		return
	}
	chatID := tgmsg.Chat.ID
	tripInput := tgmsg.Text
	msg := h.usecase.CalculateTrip(ctx, chatID, tripInput)
	if _, err := b.SendMessage(ctx, &tgapi.SendMessageParams{
		ChatID:      chatID,
		Text:        msg.Text,
		ParseMode:   tgmodel.ParseModeMarkdown,
		ReplyMarkup: msg.Markup,
	}); err != nil {
		slog.Error("SendMessage failed", "method", "defaultMessage", "chatID", chatID, "input", tripInput, "text", msg.Text, "err", err)
	}
}
