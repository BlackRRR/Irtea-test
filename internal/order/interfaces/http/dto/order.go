package dto

import "github.com/shopspring/decimal"

type OrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type PlaceOrderRequest struct {
	UserID string             `json:"user_id" validate:"required"`
	Items  []OrderItemRequest `json:"items" validate:"required,min=1"`
}

type OrderItemResponse struct {
	ProductID          string          `json:"product_id"`
	ProductDescription string          `json:"product_description"`
	ProductPrice       decimal.Decimal `json:"product_price"`
	Quantity           int             `json:"quantity"`
	TotalPrice         decimal.Decimal `json:"total_price"`
}

type OrderResponse struct {
	ID         string              `json:"id"`
	UserID     string              `json:"user_id"`
	Items      []OrderItemResponse `json:"items"`
	Status     string              `json:"status"`
	TotalPrice decimal.Decimal     `json:"total_price"`
	CreatedAt  string              `json:"created_at"`
	UpdatedAt  string              `json:"updated_at"`
}
