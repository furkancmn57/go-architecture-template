// Package apperr is the typed error currency for services and transports.
package apperr

import (
	"net/http"

	"github.com/furkancmn57/go-base-template/src/constants"
)

// Error is the typed error every service method returns instead of a raw
// Go error. It carries enough information for the transport layer to render
// a consistent JSON error response.
type Error struct {
	Status  int               `json:"-"`
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
	Err     error             `json:"-"` // wrapped cause, for logging only, never serialized
}

// Error implements the error interface so *Error can still be used with
// standard library helpers (errors.Is/As, logging, etc.) internally.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap exposes the wrapped cause for errors.Is/As support.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// New builds a plain *Error with no wrapped cause.
func New(status int, code, message string) *Error {
	return &Error{Status: status, Code: code, Message: message}
}

// Wrap builds an *Error around an existing error, keeping the original for
// logging while exposing only Code/Message to clients.
func Wrap(status int, code, message string, err error) *Error {
	return &Error{Status: status, Code: code, Message: message, Err: err}
}

// BadRequest returns a 400 error for malformed input.
func BadRequest(code, message string) *Error {
	return New(http.StatusBadRequest, code, message)
}

// Validation returns a 422 error carrying field-level validation failures.
func Validation(details map[string]string) *Error {
	return &Error{
		Status:  http.StatusUnprocessableEntity,
		Code:    constants.ValidationError,
		Message: "one or more fields are invalid",
		Details: details,
	}
}

// NotFound returns a 404 error.
func NotFound(code, message string) *Error {
	return New(http.StatusNotFound, code, message)
}

// Conflict returns a 409 error.
func Conflict(code, message string) *Error {
	return New(http.StatusConflict, code, message)
}

// Unauthorized returns a 401 error.
func Unauthorized(message string) *Error {
	return New(http.StatusUnauthorized, constants.Unauthorized, message)
}

// Forbidden returns a 403 error.
func Forbidden(message string) *Error {
	return New(http.StatusForbidden, constants.Forbidden, message)
}

// Internal returns a 500 error, wrapping the original cause for logging.
func Internal(err error) *Error {
	return Wrap(http.StatusInternalServerError, constants.InternalError, "something went wrong", err)
}
