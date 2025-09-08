package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/BlackRRR/Irtea-test/infrastructure/postgres"
	"github.com/BlackRRR/Irtea-test/internal/user/domain"
	uService "github.com/BlackRRR/Irtea-test/internal/user/app"
)

var _ uService.UserRepo = (*UserRepo)(nil)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	const query = `
		INSERT INTO users.user (id, first_name, last_name, email, age, is_married, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	userDB := FromDomain(user)

	querier := postgres.GetQuerier(ctx, r.pool)

	_, err := querier.Exec(ctx, query,
		userDB.ID,
		userDB.FirstName,
		userDB.LastName,
		userDB.Email,
		userDB.Age,
		userDB.IsMarried,
		userDB.Password,
		userDB.CreatedAt,
		userDB.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	query := `
		SELECT id, email, first_name, last_name, age, is_married, password_hash, created_at, updated_at
		FROM users."user"
		WHERE id = $1
	`

	querier := postgres.GetQuerier(ctx, r.pool)
	row := querier.QueryRow(ctx, query, id.String())

	var userDB UserDB
	err := row.Scan(
		&userDB.ID,
		&userDB.Email,
		&userDB.FirstName,
		&userDB.LastName,
		&userDB.Age,
		&userDB.IsMarried,
		&userDB.Password,
		&userDB.CreatedAt,
		&userDB.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return userDB.ToDomain()
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `
		SELECT id, first_name, last_name, age, is_married, password_hash, created_at, updated_at
		FROM users."user"
		WHERE $1
	`

	querier := postgres.GetQuerier(ctx, r.pool)
	row := querier.QueryRow(ctx, query, email)

	var userDB UserDB
	err := row.Scan(
		&userDB.ID,
		&userDB.FirstName,
		&userDB.LastName,
		&userDB.Age,
		&userDB.IsMarried,
		&userDB.Password,
		&userDB.CreatedAt,
		&userDB.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return userDB.ToDomain()
}

func (r *UserRepo) Update(ctx context.Context, user *domain.User) error {
	const query = `
		UPDATE users."user"
		SET first_name = $2, last_name = $3, age = $4, is_married = $5, password_hash = $6, updated_at = $7
		WHERE id = $1
	`

	userDB := FromDomain(user)
	querier := postgres.GetQuerier(ctx, r.pool)

	result, err := querier.Exec(ctx, query,
		userDB.ID,
		userDB.FirstName,
		userDB.LastName,
		userDB.Age,
		userDB.IsMarried,
		userDB.Password,
		userDB.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id domain.UserID) error {
	query := `DELETE FROM users."user" WHERE id = $1`

	querier := postgres.GetQuerier(ctx, r.pool)
	result, err := querier.Exec(ctx, query, id.String())

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
