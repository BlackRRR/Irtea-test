package dto

import "github.com/shopspring/decimal"

type CreateProductRequest struct {
	Description string          `json:"description" validate:"required"`
	Tags        []string        `json:"tags"`
	Price       decimal.Decimal `json:"price" validate:"required"`
	Quantity    int             `json:"quantity" validate:"required,min=0"`
}

type UpdatePriceRequest struct {
	Price decimal.Decimal `json:"price" validate:"required"`
}

type AdjustStockRequest struct {
	Quantity int `json:"quantity" validate:"required"`
}

type ProductResponse struct {
	ID          string          `json:"id"`
	Description string          `json:"description"`
	Tags        []string        `json:"tags"`
	Price       decimal.Decimal `json:"price"`
	Quantity    int             `json:"quantity"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}
