package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductID uuid.UUID

func NewProductID() ProductID {
	return ProductID(uuid.New())
}

func (id ProductID) String() string {
	return uuid.UUID(id).String()
}

type Money struct {
	amount decimal.Decimal
}

func NewMoney(amount decimal.Decimal) (Money, error) {
	if amount.IsNegative() {
		return Money{}, ErrMoneyCannotBeNeg
	}
	return Money{amount: amount}, nil
}

func (m Money) Amount() decimal.Decimal {
	return m.amount
}

func (m Money) IsZero() bool {
	return m.amount.IsZero()
}

type Inventory struct {
	quantity int
}

func NewInventory(quantity int) (Inventory, error) {
	if quantity < 0 {
		return Inventory{}, ErrInventoryQuantityCannotBeNeg
	}
	return Inventory{quantity: quantity}, nil
}

func (i *Inventory) Quantity() int {
	return i.quantity
}

func (i *Inventory) IsAvailable(requestedQuantity int) bool {
	return i.quantity >= requestedQuantity
}

func (i *Inventory) Reserve(quantity int) error {
	if !i.IsAvailable(quantity) {
		return ErrInsufficientStock
	}
	i.quantity -= quantity
	return nil
}

func (i *Inventory) Add(quantity int) error {
	if quantity <= 0 {
		return ErrQuantityToAddMustBePositive
	}
	i.quantity += quantity
	return nil
}

type Product struct {
	ID          ProductID
	Description string
	Tags        []string
	Price       Money
	Inventory   Inventory
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProduct(description string, tags []string, price Money, inventory Inventory) (*Product, error) {
	description = strings.TrimSpace(description)
	if description == "" {
		return nil, ErrProductDescCannotBeEmpty
	}

	cleanedTags := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			cleanedTags = append(cleanedTags, tag)
		}
	}

	return &Product{
		ID:          NewProductID(),
		Description: description,
		Tags:        cleanedTags,
		Price:       price,
		Inventory:   inventory,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func (p *Product) UpdatePrice(price Money) {
	p.Price = price
	p.UpdatedAt = time.Now()
}

func (p *Product) AdjustStock(quantity int) error {
	if quantity > 0 {
		return p.Inventory.Add(quantity)
	} else if quantity < 0 {
		return p.Inventory.Reserve(-quantity)
	}

	return nil
}

func (p *Product) IsAvailable(quantity int) bool {
	return p.Inventory.IsAvailable(quantity)
}
