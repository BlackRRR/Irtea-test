package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type Config struct {
	Host            string        `env:"HOST" validate:"required"`
	Port            string        `env:"PORT" validate:"required"`
	Database        string        `env:"DATABASE" validate:"required"`
	Username        string        `env:"USERNAME" validate:"required"`
	Password        string        `env:"PASSWORD"  validate:"required"`
	SSLMode         string        `env:"SSLMODE"  validate:"required"`
	MaxConnections  int32         `env:"MAX_CONNECTIONS" envDefault:"100"`
	MinConnections  int32         `env:"MIN_CONNECTIONS" envDefault:"5"`
	MaxConnLifetime time.Duration `env:"MAX_CONN_LIFETIME" env.Default:"5m"`
	MaxConnIdleTime time.Duration `env:"MAX_CONN_IDLETIME" env.Default:"5m"`
}

type DB struct {
	pool *pgxpool.Pool
}

type DBOptions struct {
	retries       int
	retryInterval time.Duration
}

func NewDBOptions() DBOptions {
	return DBOptions{
		retries:       5,
		retryInterval: time.Second * 1,
	}
}

func (o DBOptions) SetRetries(retries int) DBOptions {
	o.retries = retries

	return o
}

func (o DBOptions) SetRetryInterval(retryInterval time.Duration) DBOptions {
	o.retryInterval = retryInterval

	return o
}

func NewPgxPoolConfig(cfg Config) (*pgxpool.Config, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database, cfg.SSLMode)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	config.MaxConns = cfg.MaxConnections
	config.MinConns = cfg.MinConnections
	config.MaxConnLifetime = cfg.MaxConnLifetime
	config.MaxConnIdleTime = cfg.MaxConnIdleTime

	return config, nil
}

func Connect(ctx context.Context, logger *slog.Logger, pgCfg *pgxpool.Config, o DBOptions) (*DB, error) {
	// Логика повторных попыток подключения с проверкой через Ping
	var err error
	for i := 0; i < o.retries; i++ {
		p, err := pgxpool.NewWithConfig(ctx, pgCfg)
		if err == nil {
			// Принудительно проверяем соединение через Ping
			if pingErr := p.Ping(ctx); pingErr == nil {
				logger.InfoContext(ctx, "Successfully connected to DB",
					slog.Int("attempt", i+1),
					slog.Int("retries", o.retries))

				return &DB{pool: p}, nil // Возвращаем успешное соединение
			} else {
				err = pingErr // Если Ping не прошел, сохраняем ошибку
			}
		}

		logger.WarnContext(ctx, "Failed to connect db", slog.Int("attempt", i+1),
			slog.Int("retries", o.retries), slog.Any("error", err))

		// Задержка перед следующей попыткой
		time.Sleep(o.retryInterval)
	}

	return nil, fmt.Errorf("failed to connect to database after %d error: %v", o.retries, err)
}

func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	return db.pool.Ping(ctx)
}
