package response

import (
	"erp-2c/lib/sl"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Code    int    `json:"-"`
	Message string `json:"message,omitempty"`
	Body    any    `json:"body,omitempty"`
}

func (r Response) Render(w http.ResponseWriter, req *http.Request) error {
	render.Status(req, r.Code)
	return nil
}

func (r Response) SendResponse(w http.ResponseWriter, req *http.Request) {
	if err := render.Render(w, req, r); err != nil {
		slog.Error("failed tor render response", sl.Err(err))
	}
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

func ValidationError(err error) Response {

	var validationErr validator.ValidationErrors
	if !errors.As(err, &validationErr) {
		return Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	details := make(map[string]string)
	for _, e := range validationErr {
		field := e.Field()
		switch e.Tag() {
		case "required":
			details[field] = fmt.Sprintf("field: %s is required field", e.Field())
		case "email":
			details[field] = fmt.Sprintf("field %s is invalid format, value: %s", e.Field(), e.Value())
		case "gte":
			details[field] = fmt.Sprintf("field %s is invalid param: %s", e.Field(), e.Param())
		default:
			details[field] = fmt.Sprintf("validation error %s", e.Tag())
		}
	}
	return Response{
		Code:    http.StatusBadRequest,
		Message: "Validation failed",
		Body:    details,
	}
}
