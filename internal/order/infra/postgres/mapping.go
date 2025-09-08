package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/BlackRRR/Irtea-test/internal/order/domain"
	productDomain "github.com/BlackRRR/Irtea-test/internal/product/domain"
	userDomain "github.com/BlackRRR/Irtea-test/internal/user/domain"
)

type OrderItemDB struct {
	ID           string          `db:"id"`
	OrderID      string          `db:"order_id"`
	ProductID    string          `db:"product_id"`
	Quantity     int             `db:"quantity"`
	ProductPrice decimal.Decimal `db:"product_price"`
	CreatedAt    time.Time       `db:"created_at"`
}

type OrderDB struct {
	ID         string          `db:"id"`
	UserID     string          `db:"user_id"`
	Status     string          `db:"status"`
	TotalPrice decimal.Decimal `db:"total_price"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
}

type OrderWithItemsDB struct {
	OrderDB
	Items []OrderItemDB
}

func (o *OrderWithItemsDB) ToDomain(productDescriptions map[string]string) (*domain.Order, error) {
	id, err := uuid.Parse(o.ID)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(o.UserID)
	if err != nil {
		return nil, err
	}

	items := make([]domain.OrderItem, 0, len(o.Items))
	for _, itemDB := range o.Items {
		itemID, err := uuid.Parse(itemDB.ID)
		if err != nil {
			return nil, err
		}

		orderID, err := uuid.Parse(itemDB.OrderID)
		if err != nil {
			return nil, err
		}

		productID, err := uuid.Parse(itemDB.ProductID)
		if err != nil {
			return nil, err
		}

		price, err := productDomain.NewMoney(itemDB.ProductPrice)
		if err != nil {
			return nil, err
		}

		description := productDescriptions[itemDB.ProductID]

		item := domain.OrderItem{
			ID:                 domain.OrderItemID(itemID),
			OrderID:            domain.OrderID(orderID),
			ProductID:          productDomain.ProductID(productID),
			ProductDescription: description,
			ProductPrice:       price,
			Quantity:           itemDB.Quantity,
			CreatedAt:          itemDB.CreatedAt,
		}

		items = append(items, item)
	}

	totalPrice, err := productDomain.NewMoney(o.TotalPrice)
	if err != nil {
		return nil, err
	}

	return &domain.Order{
		ID:         domain.OrderID(id),
		UserID:     userDomain.UserID(userID),
		Items:      items,
		Status:     domain.OrderStatus(o.Status),
		TotalPrice: totalPrice,
		CreatedAt:  o.CreatedAt,
		UpdatedAt:  o.UpdatedAt,
	}, nil
}

func FromDomain(order *domain.Order) (*OrderWithItemsDB, error) {
	itemsDB := make([]OrderItemDB, 0, len(order.Items))
	for _, item := range order.Items {
		itemDB := OrderItemDB{
			ID:           item.ID.String(),
			OrderID:      item.OrderID.String(),
			ProductID:    item.ProductID.String(),
			Quantity:     item.Quantity,
			ProductPrice: item.ProductPrice.Amount(),
			CreatedAt:    item.CreatedAt,
		}
		itemsDB = append(itemsDB, itemDB)
	}

	return &OrderWithItemsDB{
		OrderDB: OrderDB{
			ID:         order.ID.String(),
			UserID:     order.UserID.String(),
			Status:     string(order.Status),
			TotalPrice: order.TotalPrice.Amount(),
			CreatedAt:  order.CreatedAt,
			UpdatedAt:  order.UpdatedAt,
		},
		Items: itemsDB,
	}, nil
}