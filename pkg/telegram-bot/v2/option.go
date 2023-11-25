package telegrambot

import (
	"cdtj.io/days-in-turkey-bot/model"
	tgapi "github.com/go-telegram/bot"
)

type TelegramBotBind struct {
	Command *model.TelegramBotCommand
	Handler tgapi.HandlerFunc
}

func NewTelegramBotBind(command *model.TelegramBotCommand, handler tgapi.HandlerFunc) *TelegramBotBind {
	return &TelegramBotBind{
		Command: command,
		Handler: handler,
	}
}

type TelegramBotOptions struct {
	BoundCommands []*TelegramBotBind
	Description   *model.TelegramBotDescription
	HndlrDefault  tgapi.HandlerFunc
	HndlrError    tgapi.ErrorsHandler
	HndlrDebug    tgapi.DebugHandler
}

func NewTelegramBotOptions(boundCommands []*TelegramBotBind, description *model.TelegramBotDescription,
	hndlrDefault tgapi.HandlerFunc, hndlrError tgapi.ErrorsHandler, hndlrDebug tgapi.DebugHandler) *TelegramBotOptions {
	return &TelegramBotOptions{
		BoundCommands: boundCommands,
		Description:   description,
		HndlrDefault:  hndlrDefault,
		HndlrError:    hndlrError,
		HndlrDebug:    hndlrDebug,
	}
}

func (o *TelegramBotOptions) GetDescription() *model.TelegramBotDescription {
	return o.Description
}

func (o *TelegramBotOptions) GetCommands() []*model.TelegramBotCommand {
	cmds := make([]*model.TelegramBotCommand, 0)
	for _, cmd := range o.BoundCommands {
		if cmd.Command != nil && cmd.Command.CommandType == model.TelegramBotCommandMessageExact && cmd.Command.Caption != "" {
			cmds = append(cmds, cmd.Command)
		}
	}
	return cmds
}

func (o *TelegramBotOptions) apiOptions() []tgapi.Option {
	opts := make([]tgapi.Option, 0)
	for _, b := range o.BoundCommands {
		opts = append(opts, b.bindHandler())
	}
	if o.HndlrDebug != nil {
		opts = append(opts, bindHandlerDefaultDebug(o.HndlrDebug))
	}
	if o.HndlrDefault != nil {
		opts = append(opts, bindHandlerDefault(o.HndlrDefault))
	}
	if o.HndlrError != nil {
		opts = append(opts, bindHandlerDefaultError(o.HndlrError))
	}
	return opts
}

func (b *TelegramBotBind) bindHandler() tgapi.Option {
	switch b.Command.CommandType {
	case model.TelegramBotCommandMessageExact:
		return bindHandlerMessageExact(b.Command.Command, b.Handler)
	case model.TelegramBotCommandMessagePrefix:
		return bindHandlerMessagePrefix(b.Command.Command, b.Handler)
	case model.TelegramBotCommandCallbackExact:
		return bindHandlerCbExact(b.Command.Command, b.Handler)
	case model.TelegramBotCommandCallbackPrefix:
		return bindHandlerCbPrefix(b.Command.Command, b.Handler)
	case model.TelegramBotCommandDefaultHandler:
		return bindHandlerDefault(b.Handler)
	default:
		return nil
	}
}

func bindHandlerMessageExact(command string, handler tgapi.HandlerFunc) tgapi.Option {
	return tgapi.WithMessageTextHandler("/"+command, tgapi.MatchTypeExact, handler)
}

func bindHandlerMessagePrefix(command string, handler tgapi.HandlerFunc) tgapi.Option {
	return tgapi.WithMessageTextHandler("/"+command, tgapi.MatchTypePrefix, handler)
}

func bindHandlerCbExact(command string, handler tgapi.HandlerFunc) tgapi.Option {
	return tgapi.WithCallbackQueryDataHandler(command, tgapi.MatchTypeExact, handler)
}

func bindHandlerCbPrefix(command string, handler tgapi.HandlerFunc) tgapi.Option {
	return tgapi.WithCallbackQueryDataHandler(command, tgapi.MatchTypePrefix, handler)
}

func bindHandlerDefault(handler tgapi.HandlerFunc) tgapi.Option {
	return tgapi.WithDefaultHandler(handler)
}

func bindHandlerDefaultError(handler tgapi.ErrorsHandler) tgapi.Option {
	return tgapi.WithErrorsHandler(handler)
}

func bindHandlerDefaultDebug(handler tgapi.DebugHandler) tgapi.Option {
	return tgapi.WithDebugHandler(handler)
}
