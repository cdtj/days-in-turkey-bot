package bot

import "errors"

var (
	ErrBotCommandNotFound = errors.New("unknown command")
	ErrBotNotImpl         = errors.New("not implemented")
)
