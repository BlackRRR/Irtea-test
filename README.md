# Irtea API

## TASK

Описание
Необходимо спроектировать и разработать сервис, который работает с реляционным
хранилищем данных PostgreSQL.

Основной упор необходимо сделать на программную архитектуру: слои, маппинги
структур данных между слоями приложения.

Основные сущности

### User

- id
- firstname
- lastname
- fullname (firstname + lastname)
- age
- is_married
- password

### Product

- id
- description
- tags
- quantity

## Что надо сделать

- Реализовать следующий функционал:
- Регистрация пользователя (не младше 18 лет);
- Пароль не меньше 8 символов;
- Пользователь может заказать продукт;
- У пользователя может быть много заказов;
- Заказ может содержать множество продуктов;
- Если, продуктов не осталось на складе – его нельзя заказать;
- Нужна историчность заказов и продуктов в заказе (например старая цена).

### Тесты

Покрыть тестами несколько (2-3) функционально важных методов.
Тезисно

**Не все из перечисленного ниже обязательно реализовывать.**

- REST API или gRPC (сделать осознанный выбор)
- Слоеная архитектура (обосновать выбор)
- Логирование (в контексте) - middleware
- (optional) Трасировка, opentelemetry - middleware
- (optional) Sentry - ловить паники в middleware
- На каждом слое своя структура данных
- Поток данных идет как в чистой или гексогональной архитектуре

# Realization

## Architecture

The application implements a layered architecture with clear separation of concerns:

Bounded Contexts:

- order,
- user,
- product

Layered architecture:

- domain (Business logic and entities)
- app Application services (use cases)
- iterfaces HTTP handlers
- infra Database repositories

## Features

### User Management

- User registration (minimum age 18, password validation)
- Authentication with bcrypt password hashing

### Product Management

- Product creation with description, tags, pricing
- Inventory management
- Stock reservation for orders

### Order Management

- Place orders with multiple items
- Order status tracking (pending → confirmed → completed/cancelled)
- Historical pricing (orders store product prices at time of purchase)
- Stock validation and reservation

## API Endpoints

### Users

- `POST /v1/users/register` - Register new user
- `GET /v1/users/{id}` - Get user by ID

### Products

- `POST /v1/products` - Create product
- `GET /v1/products` - List products (with pagination)
- `GET /v1/products/{id}` - Get product by ID
- `PUT /v1/products/{id}/price` - Update product price
- `PUT /v1/products/{id}/stock` - Adjust stock quantity

### Orders

- `POST /v1/orders` - Place new order
- `GET /v1/orders/{id}` - Get order by ID
- `GET /v1/orders/users/{userId}` - Get user's orders
- `PUT /v1/orders/{id}/confirm` - Confirm order
- `PUT /v1/orders/{id}/cancel` - Cancel order

### Health Check

- `GET /v1/health` - Service health check

## Technology Stack

- **Go 1.25** - Programming language
- **Fiber v2** - HTTP framework
- **PostgreSQL** - Database
- **pgx/v5** - PostgreSQL driver
- **OpenTelemetry** - Observability (tracing)
- **Zap/slog** - Structured logging
- **Testify** - Testing framework

## Configuration

The application uses environment variables for configuration. See `.env.example` for available options:

## Getting Started

1. **Prerequisites**:
    - Go 1.25+
    - PostgreSQL 13+

2. **Setup Database**:
   ```bash
     goose -dir ./migrations postgres "postgres://user:password@host:port/db?sslmode=disable" up       
   ```

3. **Run Application**:
   ```bash
   go run cmd/api/main.go
   ```

## Testing

Run tests for different components:

```bash
# Domain logic tests
go test ./internal/product/domain -v

# Service layer tests  
go test ./internal/user/app -v
go test ./internal/order/app -v

# All tests
go test ./... -v
```
