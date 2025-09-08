package domain

import "errors"

var (
	ErrProductNotFound              = errors.New("product not found")
	ErrInsufficientStock            = errors.New("insufficient stock")
	ErrQuantityToAddMustBePositive  = errors.New("quantity to add must be positive")
	ErrInventoryQuantityCannotBeNeg = errors.New("inventory quantity cannot be negative")
	ErrMoneyCannotBeNeg             = errors.New("money cannot be negative")
	ErrProductDescCannotBeEmpty     = errors.New("product description cannot be empty")
	ErrQuantityToAddMustBe          = errors.New("quantity must be less than zero")
	ErrInvalidPrice                 = errors.New("invalid price")
	ErrInvalidQuantity              = errors.New("invalid quantity")
)
