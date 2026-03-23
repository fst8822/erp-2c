package types

import (
	"api-gateway/lib/response"

	"google.golang.org/grpc/codes"
	status2 "google.golang.org/grpc/status"
)

func HandleError(err error) response.Response {
	status, ok := status2.FromError(err)
	if !ok {
		return response.InternalServerError()
	}

	switch status.Code() {
	case codes.NotFound:
		return response.NotFound(status.Message())
	case codes.AlreadyExists:
		return response.AlreadyExist(status.Message())
	case codes.PermissionDenied:
		return response.Forbidden(status.Message())
	case codes.Unauthenticated:
		return response.Unauthorized(status.Message())
	case codes.Internal:
		return response.InternalServerError()
	case codes.InvalidArgument:
		return response.BadRequest(status.Message())
	case codes.Aborted:
		return response.InternalServerError()
	case codes.Unknown:
		return response.InternalServerError()
	default:
		return response.InternalServerError()
	}
}
