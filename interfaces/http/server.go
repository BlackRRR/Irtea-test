package http

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	orderHandler "github.com/BlackRRR/Irtea-test/internal/order/interfaces/http"
	productHandler "github.com/BlackRRR/Irtea-test/internal/product/interfaces/http"
	userHandler "github.com/BlackRRR/Irtea-test/internal/user/interfaces/http"
	"log/slog"
	"github.com/BlackRRR/Irtea-test/interfaces/http/middleware"
)

type Config struct {
	Port         string        `env:"PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
}

type Server struct {
	app             *fiber.App
	config          Config
	logger          *slog.Logger
	middleware      *middleware.Middleware
	usersHandler    *userHandler.UsersHandler
	productsHandler *productHandler.ProductsHandler
	ordersHandler   *orderHandler.OrdersHandler
}

func NewServer(
	config Config,
	logger *slog.Logger,
	middleware *middleware.Middleware,
	usersHandler *userHandler.UsersHandler,
	productsHandler *productHandler.ProductsHandler,
	ordersHandler *orderHandler.OrdersHandler,
) *Server {
	errHandler := ErrorHandler{logger: logger}

	app := fiber.New(fiber.Config{
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		ErrorHandler: errHandler.Init()},
	)

	return &Server{
		app:             app,
		config:          config,
		logger:          logger,
		middleware:      middleware,
		usersHandler:    usersHandler,
		productsHandler: productsHandler,
		ordersHandler:   ordersHandler,
	}
}

func (s *Server) setupMiddleware() {
	s.app.Use(requestid.New())

	// Use our custom recovery middleware with Sentry integration
	s.app.Use(s.middleware.RecoveryMiddleware())

	s.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	s.app.Use(s.middleware.LoggingMiddleware())
	s.app.Use(s.middleware.TracingMiddleware())
}

func (s *Server) setupRoutes() {
	api := s.app.Group("/v1")

	api.Get("/health", s.healthCheck)

	{
		users := api.Group("/users")
		users.Post("/register", s.usersHandler.Register)
		users.Get("/:id", s.usersHandler.GetByID)
	}

	products := api.Group("/products")

	{
		products.Post("/", s.productsHandler.CreateProduct)
		products.Get("/", s.productsHandler.GetProducts)
		products.Get("/:id", s.productsHandler.GetProduct)
		products.Put("/:id/price", s.productsHandler.UpdatePrice)
		products.Put("/:id/stock", s.productsHandler.AdjustStock)
	}

	orders := api.Group("/orders")

	{
		orders.Post("/", s.ordersHandler.PlaceOrder)
		orders.Get("/:id", s.ordersHandler.GetOrder)
		orders.Put("/:id/confirm", s.ordersHandler.ConfirmOrder)
		orders.Put("/:id/cancel", s.ordersHandler.CancelOrder)
		orders.Get("/users/:userId", s.ordersHandler.GetUserOrders)
	}
}

func (s *Server) healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "irtea-api",
	})
}

func (s *Server) Start() error {
	s.setupMiddleware()
	s.setupRoutes()

	s.logger.Info("Starting HTTP server", slog.String("port", s.config.Port))

	return s.app.Listen(":" + s.config.Port)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.app.Shutdown()
}
