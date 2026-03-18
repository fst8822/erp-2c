package response

import (
	"net/http"
)

// Response todo в разработке
type Response struct {
	Code    int    `json:"-"`
	Message string `json:"message,omitempty"`
	Body    any    `json:"body,omitempty"`
}

func (r Response) SendResponse() {

}

func NoContent() Response {
	return Response{
		Code: http.StatusNoContent,
	}
}

func OK(body any) Response {
	return Response{
		Code: http.StatusOK,
		Body: &body,
	}
}

func Created(body any) Response {
	return Response{
		Code: http.StatusCreated,
		Body: &body,
	}
}

func BadRequest(message string) Response {
	return Response{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

func Unauthorized(message string) Response {
	return Response{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

func Forbidden(message string) Response {
	return Response{
		Code:    http.StatusForbidden,
		Message: message,
	}
}

func NotFound(message string) Response {
	return Response{
		Code:    http.StatusNotFound,
		Message: message,
	}
}

func AlreadyExist(message string) Response {
	return Response{
		Code:    http.StatusConflict,
		Message: message,
	}
}

func InternalServerError() Response {
	return Response{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
	}
}
