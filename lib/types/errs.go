package types

import "errors"

var (
	ErrNotFound        = errors.New("resource not found")
	ErrAlreadyExist    = errors.New("entity already exist")
	ErrForbidden       = errors.New("forbidden access")
	ErrBadRequest      = errors.New("bad request")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrInternalServer  = errors.New("internalServerError")
	ErrDatabaseTimeout = errors.New("database timeout")
	ErrNoFieldsUpdate  = errors.New("no fields to update")
)

type AppErr struct {
	Message string
	Err     error
}

func NewAppErr(message string, err error) *AppErr {
	return &AppErr{Message: message, Err: err}
}

func (a *AppErr) Error() string {
	return a.Err.Error()
}

func (a *AppErr) Unwrap() error {
	return a.Err
}
