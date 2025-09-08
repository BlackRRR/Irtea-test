package middleware

import (
	"log/slog"
	"github.com/gofiber/fiber/v2"
	"fmt"
	"github.com/BlackRRR/Irtea-test/pkg/observability/tracer"
	"time"
	"github.com/getsentry/sentry-go"
	"runtime/debug"
	sentryPkg "github.com/BlackRRR/Irtea-test/pkg/observability/sentry"
)

type Middleware struct {
	logger *slog.Logger
}

func NewMiddleware(logger *slog.Logger) *Middleware {
	return &Middleware{logger: logger}
}

func (m *Middleware) TracingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		spanName := fmt.Sprintf("%s %s", c.Method(), c.Path())
		ctx, span := tracer.Start(c.Context(), spanName,
			tracer.String("http.method", c.Method()),
			tracer.String("http.path", c.Path()),
			tracer.String("http.user_agent", c.Get("User-Agent")),
			tracer.String("http.request_id", c.Get("X-Request-Id")),
		)
		defer span.End()

		c.SetUserContext(ctx)

		return c.Next()
	}
}

func (m *Middleware) LoggingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		m.logger.Info("HTTP request",
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Duration("duration", duration),
			slog.String("user_agent", c.Get("User-Agent")),
			slog.String("request_id", c.Get("X-Request-Id")),
		)

		return err
	}
}

func (m *Middleware) RecoveryMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()

				// Log panic locally
				m.logger.Error("Panic recovered",
					slog.String("method", c.Method()),
					slog.String("path", c.Path()),
					slog.Any("panic", r),
					slog.String("stack", string(stack)),
				)

				// Send to Sentry if available
				if sentry.CurrentHub().Client() != nil {
					sentry.WithScope(func(scope *sentry.Scope) {
						scope.SetTag("method", c.Method())
						scope.SetTag("path", c.Path())
						scope.SetLevel(sentry.LevelFatal)

						sentry.CaptureException(fmt.Errorf("panic: %v", r))
					})
				}

				sentryPkg.CaptureError(fmt.Errorf("panic: %v", r), map[string]string{
					"method": c.Method(),
					"path":   c.Path(),
					"type":   "panic",
				})

				// Return 500 error
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Internal server error",
				})
			}
		}()

		return c.Next()
	}
}
