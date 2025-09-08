package logger

import (
	"log/slog"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	zapslog2 "github.com/BlackRRR/Irtea-test/pkg/observability/logger/zapslog"
	"github.com/BlackRRR/Irtea-test/pkg/environment"
)

func NewZapLogger(level LogLevel, env environment.AppEnv, logFormat LogFormat) (*slog.Logger, error) {
	logLevel, err := GetLogLevelByName(level)
	if err != nil {
		return nil, err
	}

	var encoderCfg zapcore.EncoderConfig
	var encoding string

	if logFormat == LogFormatJson {
		encoderCfg = zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoding = "json"
	} else {
		// Конфигурация для разработки (форматированный текст)
		encoderCfg = zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		encoding = "console"
	}

	development := env != environment.AppEnvProduction

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.Level(logLevel)),
		Development:      development,
		Encoding:         encoding,
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	logger := slog.New(zapslog2.NewHandler(development, zapLogger.Core(), zapslog2.WithCaller(true)))
	return logger, nil
}
