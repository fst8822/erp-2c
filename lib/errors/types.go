package types

var (
	ErrNotFound            = errors.New("resource not found")
	ErrAlreadyExist        = errors.New("Entity already exist")
	ErrForbidden           = errors.New("forbidden access")
	ErrBadRequest          = errors.New("bad request")
	ErrUnauthorized        = errors.New("unauthorized")
)

