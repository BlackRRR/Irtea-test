package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewMoney_Success(t *testing.T) {
	amount := decimal.NewFromFloat(10.50)
	money, err := NewMoney(amount)

	assert.NoError(t, err)
	assert.True(t, amount.Equal(money.Amount()))
	assert.False(t, money.IsZero())
}

func TestNewMoney_NegativeAmount(t *testing.T) {
	amount := decimal.NewFromFloat(-10.50)
	money, err := NewMoney(amount)

	assert.Error(t, err)
	assert.Equal(t, Money{}, money)
	assert.Contains(t, err.Error(), "cannot be negative")
}

func TestNewInventory_Success(t *testing.T) {
	inventory, err := NewInventory(100)

	assert.NoError(t, err)
	assert.Equal(t, 100, inventory.Quantity())
	assert.True(t, inventory.IsAvailable(50))
	assert.True(t, inventory.IsAvailable(100))
}

func TestNewInventory_NegativeQuantity(t *testing.T) {
	inventory, err := NewInventory(-10)

	assert.Error(t, err)
	assert.Equal(t, Inventory{}, inventory)
	assert.Contains(t, err.Error(), "cannot be negative")
}

func TestInventory_IsAvailable(t *testing.T) {
	inventory, _ := NewInventory(10)

	assert.True(t, inventory.IsAvailable(5))
	assert.True(t, inventory.IsAvailable(10))
	assert.False(t, inventory.IsAvailable(15))
}

func TestInventory_Reserve_Success(t *testing.T) {
	inventory, _ := NewInventory(10)

	err := inventory.Reserve(3)
	assert.NoError(t, err)
	assert.Equal(t, 7, inventory.Quantity())
}

func TestInventory_Reserve_InsufficientStock(t *testing.T) {
	inventory, _ := NewInventory(5)

	err := inventory.Reserve(10)
	assert.Error(t, err)
	assert.Equal(t, ErrInsufficientStock, err)
	assert.Equal(t, 5, inventory.Quantity()) // Quantity should remain unchanged
}

func TestInventory_Add_Success(t *testing.T) {
	inventory, _ := NewInventory(10)

	err := inventory.Add(5)
	assert.NoError(t, err)
	assert.Equal(t, 15, inventory.Quantity())
}

func TestInventory_Add_InvalidQuantity(t *testing.T) {
	inventory, _ := NewInventory(10)

	err := inventory.Add(-5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be positive")
	assert.Equal(t, 10, inventory.Quantity()) // Quantity should remain unchanged
}

func TestNewProduct_Success(t *testing.T) {
	price, _ := NewMoney(decimal.NewFromFloat(19.99))
	inventory, _ := NewInventory(50)
	tags := []string{"electronics", "gadget"}

	product, err := NewProduct("Test Product", tags, price, inventory)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "Test Product", product.Description)
	assert.Equal(t, tags, product.Tags)
	assert.True(t, price.Amount().Equal(product.Price.Amount()))
	assert.Equal(t, 50, product.Inventory.Quantity())
	assert.NotEqual(t, ProductID{}, product.ID)
}

func TestNewProduct_EmptyDescription(t *testing.T) {
	price, _ := NewMoney(decimal.NewFromFloat(19.99))
	inventory, _ := NewInventory(50)

	product, err := NewProduct("", []string{}, price, inventory)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "cannot be empty")
}

func TestNewProduct_EmptyTagsFiltered(t *testing.T) {
	price, _ := NewMoney(decimal.NewFromFloat(19.99))
	inventory, _ := NewInventory(50)
	tags := []string{"electronics", "", "  ", "gadget"}

	product, err := NewProduct("Test Product", tags, price, inventory)

	assert.NoError(t, err)
	assert.Equal(t, []string{"electronics", "gadget"}, product.Tags)
}

func TestProduct_UpdatePrice(t *testing.T) {
	price, _ := NewMoney(decimal.NewFromFloat(19.99))
	inventory, _ := NewInventory(50)
	product, _ := NewProduct("Test Product", []string{}, price, inventory)

	originalUpdatedAt := product.UpdatedAt

	newPrice, _ := NewMoney(decimal.NewFromFloat(29.99))
	product.UpdatePrice(newPrice)

	assert.True(t, newPrice.Amount().Equal(product.Price.Amount()))
	assert.True(t, product.UpdatedAt.After(originalUpdatedAt))
}

func TestProduct_AdjustStock_Positive(t *testing.T) {
	price, _ := NewMoney(decimal.NewFromFloat(19.99))
	inventory, _ := NewInventory(10)
	product, _ := NewProduct("Test Product", []string{}, price, inventory)

	err := product.AdjustStock(5)
	assert.NoError(t, err)
	assert.Equal(t, 15, product.Inventory.Quantity())
}

func TestProduct_AdjustStock_Negative(t *testing.T) {
	price, _ := NewMoney(decimal.NewFromFloat(19.99))
	inventory, _ := NewInventory(10)
	product, _ := NewProduct("Test Product", []string{}, price, inventory)

	err := product.AdjustStock(-3)
	assert.NoError(t, err)
	assert.Equal(t, 7, product.Inventory.Quantity())
}

func TestProduct_AdjustStock_InsufficientForNegative(t *testing.T) {
	price, _ := NewMoney(decimal.NewFromFloat(19.99))
	inventory, _ := NewInventory(5)
	product, _ := NewProduct("Test Product", []string{}, price, inventory)

	err := product.AdjustStock(-10)
	assert.Error(t, err)
	assert.Equal(t, ErrInsufficientStock, err)
	assert.Equal(t, 5, product.Inventory.Quantity()) // Should remain unchanged
}

func TestProduct_IsAvailable(t *testing.T) {
	price, _ := NewMoney(decimal.NewFromFloat(19.99))
	inventory, _ := NewInventory(10)
	product, _ := NewProduct("Test Product", []string{}, price, inventory)

	assert.True(t, product.IsAvailable(5))
	assert.True(t, product.IsAvailable(10))
	assert.False(t, product.IsAvailable(15))
}