package app

import (
	"context"

	"github.com/BlackRRR/Irtea-test/internal/product/domain"
)

type ProductRepo interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id domain.ProductID) (*domain.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id domain.ProductID) error
	ReserveStock(ctx context.Context, id domain.ProductID, quantity int) error
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}