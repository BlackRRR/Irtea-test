-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS products;

CREATE TABLE IF NOT EXISTS products.product
(
    id          UUID PRIMARY KEY,
    description TEXT                     NOT NULL CHECK (LENGTH(TRIM(description)) > 0),
    tags        TEXT,
    price       numeric                  NOT NULL CHECK (price >= 0),
    quantity    INTEGER                  NOT NULL CHECK (quantity >= 0),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_description ON products.product (description);
CREATE INDEX idx_products_price ON products.product (price);
CREATE INDEX idx_products_quantity ON products.product (quantity);
CREATE INDEX idx_products_created_at ON products.product (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS products CASCADE;
-- +goose StatementEnd
