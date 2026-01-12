package response

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

const (
	statusOK                  = 200
	statusCreated             = 201
	statusBadRequest          = 400
	statusUnauthorized        = 401
	statusForbidden           = 403
	statusNotFound            = 404
	statusConflict            = 409
	statusInternalServerError = 500
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Body    any    `json:"body,omitempty"`
}

func New(code int, message string, body any) *Response {
	return &Response{Code: code, Message: message, Body: body}
}

func Success(code int, msg string, body any) *Response {
	return New(code, msg, body)
}

func Error(code int, msg string) *Response {
	return New(code, msg, nil)
}

func OK(msg string, body any) *Response {
	return New(statusOK, msg, body)
}

func Created(msg string, body any) *Response {
	return New(statusCreated, msg, body)
}

func BadRequest(message string, body any) *Response {
	return New(statusBadRequest, message, body)
}

func Unauthorized(message string) *Response {
	return New(statusUnauthorized, message, nil)
}

func Forbidden(message string) *Response {
	return New(statusForbidden, message, nil)
}

func NotFound(message string) *Response {
	return New(statusNotFound, message, nil)
}

func AlreadyExist(message string) *Response {
	return New(statusConflict, message, nil)
}

func InternalServerError(message string) *Response {
	return New(statusInternalServerError, message, nil)
}

func ErrorWithData(code int, message string, data any) *Response {
	return New(code, message, data)
}

func (r *Response) WithBody(body any) {
	r.Body = body
}

func validateErrors(err validator.ValidationErrors) {
	details := make(map[string]string)

	for _, e := range err {
		field := e.Field()
		switch e.Tag() {
		case details[field]:
			fmt.Sprintf("")

		}
	}
}
