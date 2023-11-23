package model

type TelegramBotCommandRow struct {
	LanguageCode string
	Commands     []*TelegramBotCommand
}

func NewTelegramBotCommandRow(commands []*TelegramBotCommand, languageCode string) *TelegramBotCommandRow {
	return &TelegramBotCommandRow{
		LanguageCode: languageCode,
		Commands:     commands,
	}
}

type TelegramBotCommand struct {
	Caption string
	Command string
}

func NewTelegramBotCommand(caption, command string) *TelegramBotCommand {
	return &TelegramBotCommand{
		Caption: caption,
		Command: command,
	}
}

type TelegramMessage struct {
	Text   string
	Markup any
}

func NewTelegramMessage(text string, markup any) *TelegramMessage {
	return &TelegramMessage{
		Text:   text,
		Markup: markup,
	}
}

type TelegramBotDescription struct {
	Description  string
	About        string
	LanguageCode string
}

func NewTelegramBotDescription(description, about, languageCode string) *TelegramBotDescription {
	return &TelegramBotDescription{
		Description:  description,
		About:        about,
		LanguageCode: languageCode,
	}
}
