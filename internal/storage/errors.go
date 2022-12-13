package storage

import "errors"

var (
	// ErrNotFound is returned when a requested item is not found
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists is returned when a requested item already exists
	ErrAlreadyExists = errors.New("already exists")
)
