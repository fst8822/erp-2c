package types

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HandleError todo в разработке
func HandleError(err error) error {
	appErr := &AppErr{}
	if !errors.As(err, &appErr) {
		return status.Error(codes.Internal, err.Error())
	}
	switch err := appErr.Unwrap(); {
	case errors.Is(err, ErrNotFound):
		return status.Error(codes.NotFound, appErr.Message)
	case errors.Is(err, ErrAlreadyExist):
		return status.Error(codes.AlreadyExists, appErr.Message)
	case errors.Is(err, ErrPasswordHash):
		return status.Error(codes.Internal, appErr.Message)
	case errors.Is(err, ErrForbidden):
		return status.Error(codes.PermissionDenied, appErr.Message)
	case errors.Is(err, ErrBadRequest):
		return status.Error(codes.InvalidArgument, appErr.Message)
	case errors.Is(err, ErrUnauthorized):
		return status.Error(codes.Unauthenticated, appErr.Message)
	case errors.Is(err, ErrDatabaseTimeout):
		return status.Error(codes.Internal, appErr.Message)
	case errors.Is(err, ErrInspectedSQL):
		return status.Error(codes.Internal, appErr.Message)
	case errors.Is(err, ErrNoFieldsUpdate):
		return status.Error(codes.InvalidArgument, appErr.Message)
	default:
		return status.Error(codes.Internal, appErr.Message)
	}
}
