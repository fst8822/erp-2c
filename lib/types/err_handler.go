package types

import (
	"erp-2c/lib/response"
	"errors"
)

func HandleError(err error) response.Response {
	appErr := &AppErr{}
	if !errors.As(err, &appErr) {
		return response.InternalServerError()
	}
	switch err := appErr.Unwrap(); {
	case errors.Is(err, ErrNotFound):
		return response.NotFound(appErr.Message)
	case errors.Is(err, ErrAlreadyExist):
		return response.AlreadyExist(appErr.Message)
	case errors.Is(err, ErrPasswordHash):
		return response.InternalServerError()
	case errors.Is(err, ErrForbidden):
		return response.Forbidden(appErr.Message)
	case errors.Is(err, ErrBadRequest):
		return response.BadRequest(appErr.Message)
	case errors.Is(err, ErrUnauthorized):
		return response.Unauthorized(appErr.Message)
	case errors.Is(err, ErrDatabaseTimeout):
		return response.InternalServerError()
	case errors.Is(err, ErrInspectedSQL):
		return response.InternalServerError()
	case errors.Is(err, ErrNoFieldsUpdate):
		return response.BadRequest(appErr.Message)
	default:
		return response.InternalServerError()
	}
}
