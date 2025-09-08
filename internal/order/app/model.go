package app

import (
	productDomain "github.com/BlackRRR/Irtea-test/internal/product/domain"
	userDomain "github.com/BlackRRR/Irtea-test/internal/user/domain"
)

type OrderItemInput struct {
	ProductID productDomain.ProductID `json:"product_id"`
	Quantity  int                     `json:"quantity"`
}

type PlaceOrderInput struct {
	UserID userDomain.UserID `json:"user_id"`
	Items  []OrderItemInput  `json:"items"`
}
