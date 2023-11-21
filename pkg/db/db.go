package db

import "errors"

// some of database methods are only used by few entities,
// so db abstractions are defined directly on entity repo level
// to avoid dummy methods declaration

var (
	ErrDBEntryNotFound  = errors.New("entry not found")
	ErrDBBucketNotFound = errors.New("bucket not found")
	ErrDBUnknownEntity  = errors.New("unknown entity")
)
