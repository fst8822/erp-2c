package types

var (
	ErrNotFound            = errors.New("resource not found")
	ErrConflict            = errors.New("datamodel conflict")
	ErrForbidden           = errors.New("forbidden access")
	ErrBadRequest          = errors.New("bad request")
	ErrUnauthorized        = errors.New("unauthorized")
)

