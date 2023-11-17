package country

import "errors"

var (
	ErrNoFiles  = errors.New("no country files in country folder")
	ErrNotFound = errors.New("not found")
)
