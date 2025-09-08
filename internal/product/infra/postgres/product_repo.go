package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/BlackRRR/Irtea-test/infrastructure/postgres"
	"github.com/BlackRRR/Irtea-test/internal/product/domain"
	pService "github.com/BlackRRR/Irtea-test/internal/product/app"
)

var _ pService.ProductRepo = (*ProductRepo)(nil)

type ProductRepo struct {
	pool *pgxpool.Pool
}

func NewProductRepo(pool *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{pool: pool}
}

func (r *ProductRepo) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products.product (id, description, tags, price, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	productDB := FromDomain(product)
	querier := postgres.GetQuerier(ctx, r.pool)

	_, err := querier.Exec(ctx, query,
		productDB.ID,
		productDB.Description,
		productDB.Tags,
		productDB.Price,
		productDB.Quantity,
		productDB.CreatedAt,
		productDB.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (r *ProductRepo) GetByID(ctx context.Context, id domain.ProductID) (*domain.Product, error) {
	query := `
		SELECT id, description, tags, price, quantity, created_at, updated_at
		FROM products.product
		WHERE id = $1
	`

	querier := postgres.GetQuerier(ctx, r.pool)
	row := querier.QueryRow(ctx, query, id.String())

	var productDB ProductDB
	err := row.Scan(
		&productDB.ID,
		&productDB.Description,
		&productDB.Tags,
		&productDB.Price,
		&productDB.Quantity,
		&productDB.CreatedAt,
		&productDB.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}

	return productDB.ToDomain()
}

func (r *ProductRepo) GetAll(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	query := `
		SELECT id, description, tags, price, quantity, created_at, updated_at
		FROM products.product
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	querier := postgres.GetQuerier(ctx, r.pool)
	rows, err := querier.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var productDB ProductDB
		err := rows.Scan(
			&productDB.ID,
			&productDB.Description,
			&productDB.Tags,
			&productDB.Price,
			&productDB.Quantity,
			&productDB.CreatedAt,
			&productDB.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}

		product, err := productDB.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert product to domain: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return products, nil
}

func (r *ProductRepo) Update(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products.product
		SET description = $2, tags = $3, price = $4, quantity = $5, updated_at = $6
		WHERE id = $1
	`

	productDB := FromDomain(product)
	querier := postgres.GetQuerier(ctx, r.pool)

	result, err := querier.Exec(ctx, query,
		productDB.ID,
		productDB.Description,
		productDB.Tags,
		productDB.Price,
		productDB.Quantity,
		productDB.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepo) Delete(ctx context.Context, id domain.ProductID) error {
	query := `DELETE FROM products.product WHERE id = $1`

	querier := postgres.GetQuerier(ctx, r.pool)
	result, err := querier.Exec(ctx, query, id.String())

	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

func (r *ProductRepo) ReserveStock(ctx context.Context, id domain.ProductID, quantity int) error {
	query := `
		UPDATE products.product
		SET quantity = quantity - $2, updated_at = NOW()
		WHERE id = $1 AND quantity >= $2
	`

	querier := postgres.GetQuerier(ctx, r.pool)
	result, err := querier.Exec(ctx, query, id.String(), quantity)

	if err != nil {
		return fmt.Errorf("failed to reserve stock: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrInsufficientStock
	}

	return nil
}
