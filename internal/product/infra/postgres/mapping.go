package postgres

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/BlackRRR/Irtea-test/internal/product/domain"
)

type ProductDB struct {
	ID          string          `db:"id"`
	Description string          `db:"description"`
	Tags        string          `db:"tags"`
	Price       decimal.Decimal `db:"price"`
	Quantity    int             `db:"quantity"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
}

func (p *ProductDB) ToDomain() (*domain.Product, error) {
	id, err := uuid.Parse(p.ID)
	if err != nil {
		return nil, err
	}

	var tags []string
	if p.Tags != "" {
		tags = strings.Split(p.Tags, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	}

	price, err := domain.NewMoney(p.Price)
	if err != nil {
		return nil, err
	}

	inventory, err := domain.NewInventory(p.Quantity)
	if err != nil {
		return nil, err
	}

	return &domain.Product{
		ID:          domain.ProductID(id),
		Description: p.Description,
		Tags:        tags,
		Price:       price,
		Inventory:   inventory,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}

func FromDomain(product *domain.Product) *ProductDB {
	var tagsStr string
	if len(product.Tags) > 0 {
		tagsStr = strings.Join(product.Tags, ",")
	}

	return &ProductDB{
		ID:          product.ID.String(),
		Description: product.Description,
		Tags:        tagsStr,
		Price:       product.Price.Amount(),
		Quantity:    product.Inventory.Quantity(),
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}