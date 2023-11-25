package model

type TelegramBotCommandRow struct {
	commands []*TelegramBotCommand
}

func NewTelegramBotCommandRow(commands []*TelegramBotCommand) *TelegramBotCommandRow {
	return &TelegramBotCommandRow{
		commands: commands,
	}
}

func (s *TelegramBotCommandRow) GetCommands() []*TelegramBotCommand {
	return s.commands
}

type TelegramBotCommandType int

const (
	TelegramBotCommandMessageExact TelegramBotCommandType = iota
	TelegramBotCommandMessagePrefix
	TelegramBotCommandCallbackExact
	TelegramBotCommandCallbackPrefix
	TelegramBotCommandDescription
	TelegramBotCommandDefaultHandler
)

type TelegramBotCommand struct {
	caption     string
	command     string
	commandType TelegramBotCommandType
}

func NewTelegramBotCommand(caption, command string, commandType TelegramBotCommandType) *TelegramBotCommand {
	return &TelegramBotCommand{
		caption:     caption,
		command:     command,
		commandType: commandType,
	}
}

func (s *TelegramBotCommand) GetCaption() string {
	return s.caption
}
func (s *TelegramBotCommand) GetCommand() string {
	return s.command
}
func (s *TelegramBotCommand) GetCommandType() TelegramBotCommandType {
	return s.commandType
}

type TelegramMessage struct {
	text   string
	markup any
}

func NewTelegramMessage(text string, markup any) *TelegramMessage {
	return &TelegramMessage{
		text:   text,
		markup: markup,
	}
}

func (s *TelegramMessage) GetText() string {
	return s.text
}

func (s *TelegramMessage) GetMarkup() any {
	return s.markup
}

type TelegramBotDescription struct {
	description string
	about       string
}

func NewTelegramBotDescription(description, about string) *TelegramBotDescription {
	return &TelegramBotDescription{
		description: description,
		about:       about,
	}
}

func (s *TelegramBotDescription) GetDescription() string {
	return s.description
}
func (s *TelegramBotDescription) GetAbout() string {
	return s.about
}
