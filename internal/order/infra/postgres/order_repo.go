package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/BlackRRR/Irtea-test/infrastructure/postgres"
	"github.com/BlackRRR/Irtea-test/internal/order/domain"
	userDomain "github.com/BlackRRR/Irtea-test/internal/user/domain"
	"github.com/shopspring/decimal"
	"time"
	oService "github.com/BlackRRR/Irtea-test/internal/order/app"
)

var _ oService.OrderRepo = (*OrderRepo)(nil)

type OrderRepo struct {
	pool *pgxpool.Pool
}

func NewOrderRepo(pool *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{pool: pool}
}

func (r *OrderRepo) Create(ctx context.Context, order *domain.Order) error {
	orderDB, err := FromDomain(order)
	if err != nil {
		return fmt.Errorf("failed to convert order to DB format: %w", err)
	}

	q := postgres.GetQuerier(ctx, r.pool)

	orderQuery := `
		INSERT INTO orders."order" (id, user_id, status, total_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = q.Exec(ctx, orderQuery,
		orderDB.ID,
		orderDB.UserID,
		orderDB.Status,
		orderDB.TotalPrice,
		orderDB.CreatedAt,
		orderDB.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	if err = r.batchInsert(ctx, orderDB.Items, q); err != nil {
		return err
	}

	return nil
}

func (r *OrderRepo) GetByID(ctx context.Context, id domain.OrderID) (*domain.Order, error) {
	q := postgres.GetQuerier(ctx, r.pool)

	orderQuery := `
		SELECT id, user_id, status, total_price, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	var orderDB OrderDB
	err := q.QueryRow(ctx, orderQuery, id.String()).Scan(
		&orderDB.ID,
		&orderDB.UserID,
		&orderDB.Status,
		&orderDB.TotalPrice,
		&orderDB.CreatedAt,
		&orderDB.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order by ID: %w", err)
	}

	itemsQuery := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.product_price, oi.created_at,
		       p.description
		FROM orders.order_items oi
		JOIN products.products p ON oi.product_id = p.id
		WHERE oi.order_id = $1
		ORDER BY oi.created_at
	`

	rows, err := q.Query(ctx, itemsQuery, id.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()

	var items []OrderItemDB
	productDescriptions := make(map[string]string)
	for rows.Next() {
		var item OrderItemDB
		var description string

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.ProductPrice,
			&item.CreatedAt,
			&description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		items = append(items, item)
		productDescriptions[item.ProductID] = description
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	orderWithItems := &OrderWithItemsDB{
		OrderDB: orderDB,
		Items:   items,
	}

	return orderWithItems.ToDomain(productDescriptions)
}

func (r *OrderRepo) GetByUserID(ctx context.Context, userID userDomain.UserID, limit, offset int) ([]*domain.Order, error) {
	q := postgres.GetQuerier(ctx, r.pool)

	ordersQuery := `
		SELECT id, user_id, status, total_price, created_at, updated_at
		FROM orders.order
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := q.Query(ctx, ordersQuery, userID.String(), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by user ID: %w", err)
	}
	defer rows.Close()

	var orderIDs []string
	var ordersDB []OrderDB

	for rows.Next() {
		var orderDB OrderDB
		err := rows.Scan(
			&orderDB.ID,
			&orderDB.UserID,
			&orderDB.Status,
			&orderDB.TotalPrice,
			&orderDB.CreatedAt,
			&orderDB.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order row: %w", err)
		}

		ordersDB = append(ordersDB, orderDB)
		orderIDs = append(orderIDs, orderDB.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	if len(orderIDs) == 0 {
		return []*domain.Order{}, nil
	}

	itemsQuery := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.product_price, oi.created_at,
		       p.description
		FROM orders.order_items oi
		LEFT JOIN products.products p ON oi.product_id = p.id
		WHERE oi.order_id = ANY($1)
		ORDER BY oi.order_id, oi.created_at
	`

	itemRows, err := q.Query(ctx, itemsQuery, orderIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer itemRows.Close()

	orderItemsMap := make(map[string][]OrderItemDB)
	productDescriptions := make(map[string]string)

	for itemRows.Next() {
		var item OrderItemDB
		var description string

		err := itemRows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.ProductPrice,
			&item.CreatedAt,
			&description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		orderItemsMap[item.OrderID] = append(orderItemsMap[item.OrderID], item)
		productDescriptions[item.ProductID] = description
	}

	if err := itemRows.Err(); err != nil {
		return nil, fmt.Errorf("item rows iteration error: %w", err)
	}

	var orders []*domain.Order
	for _, orderDB := range ordersDB {
		orderWithItems := &OrderWithItemsDB{
			OrderDB: orderDB,
			Items:   orderItemsMap[orderDB.ID],
		}

		order, err := orderWithItems.ToDomain(productDescriptions)
		if err != nil {
			return nil, fmt.Errorf("failed to convert order to domain: %w", err)
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepo) Update(ctx context.Context, order *domain.Order) error {
	orderDB, err := FromDomain(order)
	if err != nil {
		return fmt.Errorf("failed to convert order to DB format: %w", err)
	}

	q := postgres.GetQuerier(ctx, r.pool)

	updateOrderQuery := `
		UPDATE orders.order
		SET status = $2, total_price = $3, updated_at = $4
		WHERE id = $1
	`

	result, err := q.Exec(ctx, updateOrderQuery,
		orderDB.ID,
		orderDB.Status,
		orderDB.TotalPrice,
		orderDB.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrOrderNotFound
	}

	deleteItemsQuery := `DELETE FROM orders.order_items WHERE order_id = $1`
	_, err = q.Exec(ctx, deleteItemsQuery, orderDB.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing order items: %w", err)
	}

	if err = r.batchInsert(ctx, orderDB.Items, q); err != nil {
		return err
	}

	return nil
}

func (r *OrderRepo) Delete(ctx context.Context, id domain.OrderID) error {
	query := `DELETE FROM orders.order WHERE id = $1`

	q := postgres.GetQuerier(ctx, r.pool)
	result, err := q.Exec(ctx, query, id.String())

	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrOrderNotFound
	}

	return nil
}

func (r *OrderRepo) batchInsert(ctx context.Context, orderItems []OrderItemDB, q postgres.Querier) error {
	ids := make([]string, 0, len(orderItems))
	orderIDs := make([]string, 0, len(orderItems))
	productIDs := make([]string, 0, len(orderItems))
	quantities := make([]int, 0, len(orderItems))
	prices := make([]decimal.Decimal, 0, len(orderItems))
	createdAts := make([]time.Time, 0, len(orderItems))

	for _, item := range orderItems {
		ids = append(ids, item.ID)
		orderIDs = append(orderIDs, item.OrderID)
		productIDs = append(productIDs, item.ProductID)
		quantities = append(quantities, item.Quantity)
		prices = append(prices, item.ProductPrice)
		createdAts = append(createdAts, item.CreatedAt)
	}

	query := `
	INSERT INTO orders.order_items (id, order_id, product_id, quantity, product_price, created_at)
	SELECT
		UNNEST($1::uuid[]),
		UNNEST($2::uuid[]),
		UNNEST($3::uuid[]),
		UNNEST($4::int[]),
		UNNEST($5::numeric[]),
		UNNEST($6::timestamptz[])
`

	if _, err := q.Exec(ctx, query,
		ids, orderIDs, productIDs, quantities, prices, createdAts,
	); err != nil {
		return fmt.Errorf("failed to create order items batch: %w", err)
	}

	return nil
}
