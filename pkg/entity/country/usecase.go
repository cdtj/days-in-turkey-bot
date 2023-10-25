package country

import "context"

type Usecase interface {
	Create(ctx context.Context) error
	Update(ctx context.Context) error
}
