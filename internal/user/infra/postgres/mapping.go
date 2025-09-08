package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/BlackRRR/Irtea-test/internal/user/domain"
)

type UserDB struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Age       int       `db:"age"`
	IsMarried bool      `db:"is_married"`
	Password  string    `db:"password_hash"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *UserDB) ToDomain() (*domain.User, error) {
	id, err := uuid.Parse(u.ID)
	if err != nil {
		return nil, err
	}

	fullName, err := domain.NewFullName(u.FirstName, u.LastName)
	if err != nil {
		return nil, err
	}

	passwordHash, err := domain.NewPasswordHash(u.Password)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        domain.UserID(id),
		Email:     u.Email,
		FullName:  fullName,
		Age:       u.Age,
		IsMarried: u.IsMarried,
		Password:  passwordHash,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func FromDomain(user *domain.User) *UserDB {
	return &UserDB{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FullName.FirstName,
		LastName:  user.FullName.LastName,
		Age:       user.Age,
		IsMarried: user.IsMarried,
		Password:  user.Password.Value(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
