package apperr

import (
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type httpBody struct {
	Error *Error `json:"error"`
}

// WriteHTTP renders an *apperr.Error as the standard JSON error envelope.
// This is the ONLY place allowed to translate an *apperr.Error into an HTTP
// response; HTTP entry points must call this instead of building error JSON
// themselves.
func WriteHTTP(c *fiber.Ctx, err *Error) error {
	if err == nil {
		err = Internal(nil)
	}
	if err.Status == 0 {
		err.Status = http.StatusInternalServerError
	}
	if err.Err != nil {
		slog.Error("apperr", "message", err.Message, "code", err.Code, "status", err.Status, "error", err.Err)
	}
	return c.Status(err.Status).JSON(httpBody{Error: err})
}
