package app

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/BlackRRR/Irtea-test/internal/order/domain"
	productDomain "github.com/BlackRRR/Irtea-test/internal/product/domain"
	userDomain "github.com/BlackRRR/Irtea-test/internal/user/domain"
)

type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) Create(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepo) GetByID(ctx context.Context, id domain.OrderID) (*domain.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepo) GetByUserID(ctx context.Context, userID userDomain.UserID, limit, offset int) ([]*domain.Order, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *MockOrderRepo) Update(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepo) Delete(ctx context.Context, id domain.OrderID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) GetByID(ctx context.Context, id productDomain.ProductID) (*productDomain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*productDomain.Product), args.Error(1)
}

func (m *MockProductRepo) ReserveStock(ctx context.Context, id productDomain.ProductID, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

type MockOrderTxManager struct {
	mock.Mock
}

func (m *MockOrderTxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	if args.Get(0) != nil {
		return args.Get(0).(func(context.Context, func(context.Context) error) error)(ctx, fn)
	}
	return fn(ctx)
}

func TestOrderService_PlaceOrder_Success(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	mockTx := new(MockOrderTxManager)

	service := NewOrderService(mockOrderRepo, mockProductRepo, mockTx)

	userID := userDomain.NewUserID()
	productID := productDomain.NewProductID()

	price, _ := productDomain.NewMoney(decimal.NewFromFloat(10.50))
	inventory, _ := productDomain.NewInventory(100)
	product, _ := productDomain.NewProduct("Test Product", []string{"tag1"}, price, inventory)
	product.ID = productID

	input := PlaceOrderInput{
		UserID: userID,
		Items: []OrderItemInput{
			{
				ProductID: productID,
				Quantity:  2,
			},
		},
	}

	mockTx.On("WithTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockProductRepo.On("GetByID", mock.Anything, productID).Return(product, nil)
	mockProductRepo.On("ReserveStock", mock.Anything, productID, 2).Return(nil)
	mockOrderRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Order")).Return(nil)

	order, err := service.PlaceOrder(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, userID, order.UserID)
	assert.Len(t, order.Items, 1)
	assert.Equal(t, productID, order.Items[0].ProductID)
	assert.Equal(t, 2, order.Items[0].Quantity)
	assert.Equal(t, domain.OrderStatusPending, order.Status)

	expectedTotal := decimal.NewFromFloat(21.00) // 10.50 * 2
	assert.True(t, expectedTotal.Equal(order.TotalPrice.Amount()))

	mockOrderRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestOrderService_PlaceOrder_InsufficientStock(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	mockTx := new(MockOrderTxManager)

	service := NewOrderService(mockOrderRepo, mockProductRepo, mockTx)

	userID := userDomain.NewUserID()
	productID := productDomain.NewProductID()

	price, _ := productDomain.NewMoney(decimal.NewFromFloat(10.50))
	inventory, _ := productDomain.NewInventory(1) // Only 1 in stock
	product, _ := productDomain.NewProduct("Test Product", []string{"tag1"}, price, inventory)
	product.ID = productID

	input := PlaceOrderInput{
		UserID: userID,
		Items: []OrderItemInput{
			{
				ProductID: productID,
				Quantity:  5, // Requesting more than available
			},
		},
	}

	mockTx.On("WithTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockProductRepo.On("GetByID", mock.Anything, productID).Return(product, nil)

	order, err := service.PlaceOrder(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, productDomain.ErrInsufficientStock, err)

	mockProductRepo.AssertExpectations(t)
	mockOrderRepo.AssertNotCalled(t, "Create")
}

func TestOrderService_PlaceOrder_ProductNotFound(t *testing.T) {
	mockOrderRepo := new(MockOrderRepo)
	mockProductRepo := new(MockProductRepo)
	mockTx := new(MockOrderTxManager)

	service := NewOrderService(mockOrderRepo, mockProductRepo, mockTx)

	userID := userDomain.NewUserID()
	productID := productDomain.NewProductID()

	input := PlaceOrderInput{
		UserID: userID,
		Items: []OrderItemInput{
			{
				ProductID: productID,
				Quantity:  1,
			},
		},
	}

	mockTx.On("WithTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockProductRepo.On("GetByID", mock.Anything, productID).Return(nil, productDomain.ErrProductNotFound)

	order, err := service.PlaceOrder(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, order)
	assert.Equal(t, productDomain.ErrProductNotFound, err)

	mockProductRepo.AssertExpectations(t)
	mockOrderRepo.AssertNotCalled(t, "Create")
}