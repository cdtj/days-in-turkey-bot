package model

type TelegramBotCommandRow struct {
	Commands []*TelegramBotCommand
}

func NewTelegramBotCommandRow(commands []*TelegramBotCommand) *TelegramBotCommandRow {
	return &TelegramBotCommandRow{
		Commands: commands,
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
	Text string
}

func NewTelegramMessage(text string) *TelegramMessage {
	return &TelegramMessage{
		Text: text,
	}
}
