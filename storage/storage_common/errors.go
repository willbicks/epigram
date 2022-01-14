package storagecommon

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrAllreadyExists = errors.New("allready exists")
)
