package errsport

import "errors"

var (
	ErrFileAlreadyExists = errors.New("file already exists")
	ErrConflictData      = errors.New("conflict data")
	ErrNotFoundData      = errors.New("not found data")
)
