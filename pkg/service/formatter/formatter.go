package formatter

import (
	"cdtj.io/days-in-turkey-bot/model"
)

type Formatter interface {
	TripTree(tree *model.TripTree) string
	User(user *model.User) string
}
