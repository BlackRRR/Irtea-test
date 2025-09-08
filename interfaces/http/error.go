package http

import (
	"github.com/gofiber/fiber/v2"
	"errors"
	"log/slog"
)

type ErrorHandler struct {
	logger *slog.Logger
}

func (e *ErrorHandler) Init() func(ctx *fiber.Ctx, err error) error {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		var errw *fiber.Error
		if errors.As(err, &errw) {
			code = errw.Code
		}

		e.logger.Error("HTTP error",
			slog.Any("error", err),
			slog.Int("status_code", code),
			slog.String("path", ctx.Path()),
			slog.String("method", ctx.Method()),
		)

		return ctx.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
