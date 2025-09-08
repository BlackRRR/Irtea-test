package app

import (
	"github.com/shopspring/decimal"
	"github.com/BlackRRR/Irtea-test/internal/product/domain"
)

type CreateProductInput struct {
	Description string          `json:"description"`
	Tags        []string        `json:"tags"`
	Price       decimal.Decimal `json:"price"`
	Quantity    int             `json:"quantity"`
}

type UpdatePriceInput struct {
	ProductID domain.ProductID `json:"product_id"`
	Price     decimal.Decimal  `json:"price"`
}

type AdjustStockInput struct {
	ProductID domain.ProductID `json:"product_id"`
	Quantity  int              `json:"quantity"`
}
