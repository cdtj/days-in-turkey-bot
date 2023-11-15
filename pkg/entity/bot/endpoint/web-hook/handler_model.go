package webhook

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
