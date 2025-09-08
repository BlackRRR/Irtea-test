package app

import (
	"context"

	"github.com/BlackRRR/Irtea-test/internal/order/domain"
	productDomain "github.com/BlackRRR/Irtea-test/internal/product/domain"
	userDomain "github.com/BlackRRR/Irtea-test/internal/user/domain"
)

type OrderRepo interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id domain.OrderID) (*domain.Order, error)
	GetByUserID(ctx context.Context, userID userDomain.UserID, limit, offset int) ([]*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) error
	Delete(ctx context.Context, id domain.OrderID) error
}

type ProductRepo interface {
	GetByID(ctx context.Context, id productDomain.ProductID) (*productDomain.Product, error)
	ReserveStock(ctx context.Context, id productDomain.ProductID, quantity int) error
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
