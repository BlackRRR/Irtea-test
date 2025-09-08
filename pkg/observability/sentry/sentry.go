package sentry

import (
	"log/slog"
	"time"

	"github.com/getsentry/sentry-go"
)

// InitSentry инициализирует Sentry если предоставлен DSN
func InitSentry(dsn, appName, environment string, logger *slog.Logger) error {
	if dsn == "" {
		logger.Info("Sentry DSN not provided, skipping Sentry initialization")
		return nil
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      environment,
		ServerName:       appName,
		TracesSampleRate: 1.0,
		AttachStacktrace: true,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// Можно добавить кастомную логику фильтрации событий
			return event
		},
	})

	if err != nil {
		logger.Error("Failed to initialize Sentry", slog.String("error", err.Error()))
		return err
	}

	logger.Info("Sentry initialized successfully",
		slog.String("environment", environment),
		slog.String("app_name", appName),
	)

	return nil
}

// Close правильно завершает работу Sentry
func Close() {
	sentry.Flush(time.Second * 2)
}

// CaptureError отправляет ошибку в Sentry (если инициализирован)
func CaptureError(err error, tags map[string]string) {
	if sentry.CurrentHub().Client() == nil {
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		for key, value := range tags {
			scope.SetTag(key, value)
		}
		sentry.CaptureException(err)
	})
}

// CaptureMessage отправляет сообщение в Sentry (если инициализирован)
func CaptureMessage(message string, level sentry.Level, tags map[string]string) {
	if sentry.CurrentHub().Client() == nil {
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(level)
		for key, value := range tags {
			scope.SetTag(key, value)
		}
		sentry.CaptureMessage(message)
	})
}
