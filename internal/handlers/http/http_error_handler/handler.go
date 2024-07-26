// Package http_error_handler catch any http errors
package http_error_handler

import (
	"github.com/gofiber/fiber/v3"

	"errors"
)

// ErrorsHandler acts like global error handler
type ErrorsHandler struct {
}

// NewErrorsHandler creates new handler
func NewErrorsHandler() *ErrorsHandler {
	return &ErrorsHandler{}
}

// Handle errors
func (hch *ErrorsHandler) Handle(ctx fiber.Ctx, err error) error {
	sc := fiber.StatusInternalServerError

	var e *fiber.Error

	if errors.As(err, &e) {
		sc = e.Code
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	ctx.Status(sc)

	return ctx.Send([]byte{})
}
