package app

import (
	"github.com/BlackRRR/Irtea-test/pkg/environment"
	"github.com/BlackRRR/Irtea-test/pkg/observability/logger"
	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/BlackRRR/Irtea-test/infrastructure/postgres"
	"github.com/BlackRRR/Irtea-test/interfaces/http"
)

type Config struct {
	AppName string `env:"APP_NAME" envDefault:"IrteaTest"`

	HttpServer http.Config `envPrefix:"HTTP_SERVER_" validate:"required"`

	// App mode: local, develop, test, production
	AppEnv environment.AppEnv `env:"APP_ENV" envDefault:"develop" validate:"required"`
	// values: debug, info, warn, error
	LogLevel logger.LogLevel `env:"LOG_LEVEL" envDefault:"debug" validate:"required"`
	// values: json, console
	LogFormat logger.LogFormat `env:"LOG_FORMAT" envDefault:"json" validate:"required"`

	Postgres postgres.Config `envPrefix:"DB_CONFIG_"`

	OtelURL string `env:"OTEL_URL"`

	// Sentry DSN (optional)
	SentryDSN string `env:"SENTRY_DSN"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	if err := validator.New().Struct(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
