package apperror

import (
	"errors"
	"net/http"
)

type AppError struct {
	HTTPStatus int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    any    `json:"details,omitempty"`
	cause      error
}

func (e *AppError) Error() string {
	if e.cause != nil {
		return e.cause.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.cause }

func (e *AppError) WithCause(err error) *AppError {
	cp := *e
	cp.cause = err
	return &cp
}

func (e *AppError) WithDetails(details any) *AppError {
	cp := *e
	cp.Details = details
	return &cp
}

func New(status int, code, message string) *AppError {
	return &AppError{HTTPStatus: status, Code: code, Message: message}
}

func NotFound(resource string) *AppError {
	return &AppError{
		HTTPStatus: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    resource + " not found",
	}
}

func BadRequest(message string) *AppError {
	return &AppError{
		HTTPStatus: http.StatusBadRequest,
		Code:       "BAD_REQUEST",
		Message:    message,
	}
}

func Internal(cause error) *AppError {
	return &AppError{
		HTTPStatus: http.StatusInternalServerError,
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "an unexpected error occurred",
		cause:      cause,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		HTTPStatus: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    message,
	}
}

func As(err error) *AppError {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae
	}
	return nil
}
