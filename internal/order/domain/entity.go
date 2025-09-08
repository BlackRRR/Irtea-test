package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	productDomain "github.com/BlackRRR/Irtea-test/internal/product/domain"
	userDomain "github.com/BlackRRR/Irtea-test/internal/user/domain"
)

type OrderItemID uuid.UUID

func NewOrderItemID() OrderItemID {
	return OrderItemID(uuid.New())
}

func (id OrderItemID) String() string {
	return uuid.UUID(id).String()
}

type OrderID uuid.UUID

func NewOrderID() OrderID {
	return OrderID(uuid.New())
}

func (id OrderID) String() string {
	return uuid.UUID(id).String()
}

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusCompleted OrderStatus = "completed"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusPending, OrderStatusConfirmed, OrderStatusCancelled, OrderStatusCompleted:
		return true
	default:
		return false
	}
}

type OrderItem struct {
	ID                 OrderItemID
	OrderID            OrderID
	ProductID          productDomain.ProductID
	ProductDescription string
	ProductPrice       productDomain.Money
	Quantity           int
	CreatedAt          time.Time
}

func NewOrderItem(orderID OrderID, productID productDomain.ProductID, description string, price productDomain.Money, quantity int) (*OrderItem, error) {
	if quantity <= 0 {
		return nil, errors.New("order item quantity must be positive")
	}

	return &OrderItem{
		ID:                 NewOrderItemID(),
		OrderID:            orderID,
		ProductID:          productID,
		ProductDescription: description,
		ProductPrice:       price,
		Quantity:           quantity,
		CreatedAt:          time.Now(),
	}, nil
}

func (oi *OrderItem) TotalPrice() productDomain.Money {
	totalAmount := oi.ProductPrice.Amount().Mul(decimal.NewFromInt(int64(oi.Quantity)))
	money, _ := productDomain.NewMoney(totalAmount)
	return money
}

type Order struct {
	ID         OrderID
	UserID     userDomain.UserID
	Items      []OrderItem
	Status     OrderStatus
	TotalPrice productDomain.Money
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewOrder(userID userDomain.UserID, items []OrderItem) (*Order, error) {
	if len(items) == 0 {
		return nil, ErrEmptyOrder
	}

	totalAmount := decimal.Zero
	for _, item := range items {
		itemTotal := item.TotalPrice()
		totalAmount = totalAmount.Add(itemTotal.Amount())
	}

	totalPrice, err := productDomain.NewMoney(totalAmount)
	if err != nil {
		return nil, err
	}

	return &Order{
		ID:         NewOrderID(),
		UserID:     userID,
		Items:      items,
		Status:     OrderStatusPending,
		TotalPrice: totalPrice,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

func (o *Order) Confirm() error {
	if o.Status != OrderStatusPending {
		return ErrInvalidOrderStatus
	}
	o.Status = OrderStatusConfirmed
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) Cancel() error {
	if o.Status == OrderStatusCompleted || o.Status == OrderStatusCancelled {
		return ErrInvalidOrderStatus
	}
	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) Complete() error {
	if o.Status != OrderStatusConfirmed {
		return ErrInvalidOrderStatus
	}
	o.Status = OrderStatusCompleted
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) CanBeModified() bool {
	return o.Status == OrderStatusPending
}