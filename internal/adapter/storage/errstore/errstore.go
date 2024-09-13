package errstore

import "errors"

var (
	ErrLoginNotUnique  = errors.New("login not unique")
	ErrNotFoundData    = errors.New("not found data")
	ErrConflictData    = errors.New("conflict data")
	ErrForbidden       = errors.New("forbidden")
	ErrInvalidImageExt = errors.New("invalid image extension")
)
