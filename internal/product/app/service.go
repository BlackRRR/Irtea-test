package app

import (
	"context"

	"github.com/BlackRRR/Irtea-test/internal/product/domain"
)

type ProductService struct {
	productRepo ProductRepo
	txManager   TxManager
}

func NewProductService(productRepo ProductRepo, txManager TxManager) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		txManager:   txManager,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, input CreateProductInput) (*domain.Product, error) {
	price, err := domain.NewMoney(input.Price)
	if err != nil {
		return nil, err
	}

	inventory, err := domain.NewInventory(input.Quantity)
	if err != nil {
		return nil, err
	}

	product, err := domain.NewProduct(input.Description, input.Tags, price, inventory)
	if err != nil {
		return nil, err
	}

	err = s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		return s.productRepo.Create(txCtx, product)
	})

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id domain.ProductID) (*domain.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *ProductService) GetProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	return s.productRepo.GetAll(ctx, limit, offset)
}

func (s *ProductService) UpdatePrice(ctx context.Context, input UpdatePriceInput) (*domain.Product, error) {
	price, err := domain.NewMoney(input.Price)
	if err != nil {
		return nil, err
	}

	var updatedProduct *domain.Product
	err = s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		product, err := s.productRepo.GetByID(txCtx, input.ProductID)
		if err != nil {
			return err
		}

		product.UpdatePrice(price)

		err = s.productRepo.Update(txCtx, product)
		if err != nil {
			return err
		}

		updatedProduct = product
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (s *ProductService) AdjustStock(ctx context.Context, input AdjustStockInput) (*domain.Product, error) {
	var updatedProduct *domain.Product
	err := s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		product, err := s.productRepo.GetByID(txCtx, input.ProductID)
		if err != nil {
			return err
		}

		err = product.AdjustStock(input.Quantity)
		if err != nil {
			return err
		}

		err = s.productRepo.Update(txCtx, product)
		if err != nil {
			return err
		}

		updatedProduct = product
		return nil
	})

	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (s *ProductService) ReserveStock(ctx context.Context, productID domain.ProductID, quantity int) error {
	return s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		return s.productRepo.ReserveStock(txCtx, productID, quantity)
	})
}
