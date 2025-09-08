package logger

import (
	"fmt"

	"github.com/pkg/errors"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type Level int8

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
)

var (
	logLevelsMap = map[LogLevel]Level{
		LogLevelDebug: DebugLevel,
		LogLevelInfo:  InfoLevel,
		LogLevelWarn:  WarnLevel,
		LogLevelError: ErrorLevel,
	}
)

func GetLogLevelByName(name LogLevel) (Level, error) {
	level, ok := logLevelsMap[name]
	if !ok {
		return DebugLevel, errors.New(fmt.Sprintf("log level `%s` not found, set `debug` level", name))
	}

	return level, nil
}
