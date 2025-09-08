package domain

import "errors"

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrEmptyOrder          = errors.New("order cannot be empty")
	ErrInvalidOrderStatus  = errors.New("invalid order status transition")
	ErrOrderCannotBeModified = errors.New("order cannot be modified in current status")
)