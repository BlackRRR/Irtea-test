package app

import (
	"context"
	obs "github.com/BlackRRR/Irtea-test/pkg/observability/logger"
	"log"
	"github.com/BlackRRR/Irtea-test/interfaces/http"
	uHandler "github.com/BlackRRR/Irtea-test/internal/user/interfaces/http"
	"github.com/BlackRRR/Irtea-test/infrastructure/postgres"
	uRepo "github.com/BlackRRR/Irtea-test/internal/user/infra/postgres"
	"github.com/BlackRRR/Irtea-test/internal/user/infra/security"
	pService "github.com/BlackRRR/Irtea-test/internal/product/app"
	pRepo "github.com/BlackRRR/Irtea-test/internal/product/infra/postgres"
	pHandler "github.com/BlackRRR/Irtea-test/internal/product/interfaces/http"

	oService "github.com/BlackRRR/Irtea-test/internal/order/app"
	oRepo "github.com/BlackRRR/Irtea-test/internal/order/infra/postgres"
	oHandler "github.com/BlackRRR/Irtea-test/internal/order/interfaces/http"
	uService "github.com/BlackRRR/Irtea-test/internal/user/app"
	"os/signal"
	"syscall"
	"log/slog"
	"github.com/BlackRRR/Irtea-test/pkg/observability/tracer"
	"time"
	"github.com/BlackRRR/Irtea-test/interfaces/http/middleware"
	sentryPkg "github.com/BlackRRR/Irtea-test/pkg/observability/sentry"
	"github.com/BlackRRR/Irtea-test/pkg/validator"
)

type App struct {
	http   *http.Server
	logger *slog.Logger
}

func InternalInit() {
	validator.Init()
}

func Init(ctx context.Context, cfg *Config) App {
	logger, err := obs.NewZapLogger(cfg.LogLevel, cfg.AppEnv, cfg.LogFormat)
	if err != nil {
		log.Fatal(err)
	}

	pgxConfig, err := postgres.NewPgxPoolConfig(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	dbOpts := postgres.NewDBOptions().SetRetryInterval(time.Second * 2)

	db, err := postgres.Connect(ctx, logger, pgxConfig, dbOpts)
	if err != nil {
		log.Fatal(err)
	}

	if err := sentryPkg.InitSentry(cfg.SentryDSN, cfg.AppName, string(cfg.AppEnv), logger); err != nil {
		logger.WarnContext(ctx, "Failed to initialize Sentry", slog.Any("error", err))
	}

	tp := tracer.InitTracer(cfg.AppName, cfg.OtelURL)
	defer func() {
		if tp != nil {
			if err = tp.Shutdown(ctx); err != nil {
				logger.ErrorContext(ctx, "Error shutting down tracer provider", slog.Any("error", err))
			}
		}
		// Close Sentry
		sentryPkg.Close()
	}()

	txManager := postgres.NewTxManager(db.Pool())

	// user
	userRepo := uRepo.NewUserRepo(db.Pool())
	userService := uService.NewUserService(userRepo, security.NewPasswordHasher(), txManager)
	userHandler := uHandler.NewUsersHandler(userService)

	// Product
	productRepo := pRepo.NewProductRepo(db.Pool())
	productService := pService.NewProductService(productRepo, txManager)
	productHandler := pHandler.NewProductsHandler(productService)

	// order
	orderRepo := oRepo.NewOrderRepo(db.Pool())
	orderService := oService.NewOrderService(orderRepo, productRepo, txManager)
	orderHandler := oHandler.NewOrdersHandler(orderService)

	mw := middleware.NewMiddleware(logger)

	server := http.NewServer(cfg.HttpServer, logger, mw, userHandler, productHandler, orderHandler)

	return App{http: server, logger: logger}
}

func (a App) Run(ctx context.Context) {
	// Channel to listen for interrupt signal to terminate server
	donec := make(chan error, 1)
	ctx, cancel := signal.NotifyContext(
		ctx,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
	defer cancel()

	// Start server in goroutine
	go func() {
		if err := a.http.Start(); err != nil {
			a.logger.ErrorContext(ctx, "Error starting server:", slog.Any("error", err))
			return
		}
	}()

	a.logger.InfoContext(ctx, "Server started successfully")

	// Graceful shutdown.
	go func() {
		select {
		case <-ctx.Done():
			donec <- ctx.Err()
		}
	}()

	<-donec

	// Block until we receive our signal
	a.logger.InfoContext(ctx, "Shutting down server...")

	// Shutdown server gracefully
	if err := a.http.Shutdown(ctx); err != nil {
		a.logger.ErrorContext(ctx, "Error shutting down server", slog.Any("error", err))
		return
	}

	a.logger.InfoContext(ctx, "Server stopped")
}
