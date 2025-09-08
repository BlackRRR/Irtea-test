-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA orders;

CREATE TYPE orders.status AS ENUM ('pending', 'confirmed', 'cancelled', 'completed');

CREATE TABLE IF NOT EXISTS orders.order
(
    id          UUID PRIMARY KEY,
    user_id     UUID                     NOT NULL,
    status      orders.status            NOT NULL,
    total_price numeric                  NOT NULL CHECK (total_price >= 0),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_orders_user_id FOREIGN KEY (user_id) REFERENCES users.user (id) ON DELETE CASCADE
);

CREATE INDEX idx_orders_user_id ON orders.order (user_id);
CREATE INDEX idx_orders_status ON orders.order (status);
CREATE INDEX idx_orders_created_at ON orders.order (created_at);


CREATE TABLE IF NOT EXISTS orders.order_items
(
    id            UUID PRIMARY KEY                  DEFAULT gen_random_uuid(),
    order_id      UUID                     NOT NULL,
    product_id    UUID                     NOT NULL,
    quantity      INTEGER                  NOT NULL CHECK (quantity > 0),
    product_price numeric                  NOT NULL CHECK (product_price >= 0),
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_order_items_order_id FOREIGN KEY (order_id) REFERENCES orders.order (id) ON DELETE CASCADE,
    CONSTRAINT fk_order_items_product_id FOREIGN KEY (product_id) REFERENCES products.product (id) ON DELETE RESTRICT
);

-- Create indexes for faster queries
CREATE INDEX idx_order_items_order_id ON orders.order_items (order_id);
CREATE INDEX idx_order_items_product_id ON orders.order_items (product_id);
CREATE INDEX idx_order_items_created_at ON orders.order_items (created_at);

CREATE UNIQUE INDEX idx_order_items_order_product ON orders.order_items (order_id, product_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS orders CASCADE;
-- +goose StatementEnd
