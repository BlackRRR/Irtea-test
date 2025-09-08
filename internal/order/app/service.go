package app

import (
	"context"

	"github.com/BlackRRR/Irtea-test/internal/order/domain"
	productDomain "github.com/BlackRRR/Irtea-test/internal/product/domain"
	userDomain "github.com/BlackRRR/Irtea-test/internal/user/domain"
)

type OrderService struct {
	orderRepo   OrderRepo
	productRepo ProductRepo
	txManager   TxManager
}

func NewOrderService(orderRepo OrderRepo, productRepo ProductRepo, txManager TxManager) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		txManager:   txManager,
	}
}

func (s *OrderService) PlaceOrder(ctx context.Context, input PlaceOrderInput) (*domain.Order, error) {
	var createdOrder *domain.Order

	err := s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		orderID := domain.NewOrderID()
		orderItems := make([]domain.OrderItem, 0, len(input.Items))

		for _, itemInput := range input.Items {
			product, err := s.productRepo.GetByID(txCtx, itemInput.ProductID)
			if err != nil {
				return err
			}

			if !product.IsAvailable(itemInput.Quantity) {
				return productDomain.ErrInsufficientStock
			}

			orderItem, err := domain.NewOrderItem(
				orderID,
				itemInput.ProductID,
				product.Description,
				product.Price,
				itemInput.Quantity,
			)
			if err != nil {
				return err
			}

			orderItems = append(orderItems, *orderItem)

			err = s.productRepo.ReserveStock(txCtx, itemInput.ProductID, itemInput.Quantity)
			if err != nil {
				return err
			}
		}

		order, err := domain.NewOrder(input.UserID, orderItems)
		if err != nil {
			return err
		}

		order.ID = orderID

		err = s.orderRepo.Create(txCtx, order)
		if err != nil {
			return err
		}

		createdOrder = order
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id domain.OrderID) (*domain.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID userDomain.UserID, limit, offset int) ([]*domain.Order, error) {
	return s.orderRepo.GetByUserID(ctx, userID, limit, offset)
}

func (s *OrderService) ConfirmOrder(ctx context.Context, id domain.OrderID) (*domain.Order, error) {
	var updatedOrder *domain.Order

	err := s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		order, err := s.orderRepo.GetByID(txCtx, id)
		if err != nil {
			return err
		}

		err = order.Confirm()
		if err != nil {
			return err
		}

		err = s.orderRepo.Update(txCtx, order)
		if err != nil {
			return err
		}

		updatedOrder = order
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, id domain.OrderID) (*domain.Order, error) {
	var cancelledOrder *domain.Order

	err := s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		order, err := s.orderRepo.GetByID(txCtx, id)
		if err != nil {
			return err
		}

		err = order.Cancel()
		if err != nil {
			return err
		}

		err = s.orderRepo.Update(txCtx, order)
		if err != nil {
			return err
		}

		cancelledOrder = order
		return nil
	})

	if err != nil {
		return nil, err
	}

	return cancelledOrder, nil
}
