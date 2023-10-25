package user

import "context"

type Service interface {
	Calc(ctx context.Context, input string, daysLimit, daysCont, resetInterval int) (string, error)
}
