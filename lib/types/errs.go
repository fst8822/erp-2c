package types

import "errors"

var (
	ErrNotFound        = errors.New("resource not found")
	ErrPasswordHash    = errors.New("failed generate password hash")
	ErrGeneratedToken  = errors.New("failed generate token jwt")
	ErrAlreadyExist    = errors.New("entity already exist")
	ErrForbidden       = errors.New("forbidden access")
	ErrBadRequest      = errors.New("bad request")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrDatabaseTimeout = errors.New("database timeout")
	ErrInspectedSQL    = errors.New("inspected SQL")
	ErrNoFieldsUpdate  = errors.New("no fields to update")
	ErrUnknownStatus   = errors.New("unknown delivery status")
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
