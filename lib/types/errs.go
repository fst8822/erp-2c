package types

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrAlreadyExist = errors.New("entity already exist")
	ErrForbidden    = errors.New("forbidden access")
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
)
