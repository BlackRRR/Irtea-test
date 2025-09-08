-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE IF NOT EXISTS users.user
(
    id            UUID PRIMARY KEY,
    first_name    VARCHAR(100)             NOT NULL,
    last_name     VARCHAR(100)             NOT NULL,
    age           INTEGER                  NOT NULL CHECK (age >= 18),
    email         varchar(100)             NOT NULL,
    is_married    BOOLEAN                  NOT NULL DEFAULT false,
    password_hash VARCHAR(255)             NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_email UNIQUE (email)
);

CREATE INDEX idx_users_email ON users.user (email);
CREATE INDEX idx_users_created_at ON users.user (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS users CASCADE;
-- +goose StatementEnd
