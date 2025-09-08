package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserID uuid.UUID

func NewUserID() UserID {
	return UserID(uuid.New())
}

func (id UserID) String() string {
	return uuid.UUID(id).String()
}

type FullName struct {
	FirstName string
	LastName  string
}

func NewFullName(firstName, lastName string) (FullName, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	if firstName == "" {
		return FullName{}, errors.New("first name cannot be empty")
	}
	if lastName == "" {
		return FullName{}, errors.New("last name cannot be empty")
	}

	return FullName{
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}

func (fn FullName) String() string {
	return fmt.Sprintf("%s %s", fn.FirstName, fn.LastName)
}

type PasswordHash struct {
	value string
}

func NewPasswordHash(hashedPassword string) (PasswordHash, error) {
	if hashedPassword == "" {
		return PasswordHash{}, errors.New("password hash cannot be empty")
	}
	return PasswordHash{value: hashedPassword}, nil
}

func (ph PasswordHash) Value() string {
	return ph.value
}

type User struct {
	ID        UserID
	Email     string
	FullName  FullName
	Age       int
	IsMarried bool
	Password  PasswordHash
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(fullName FullName, age int, isMarried bool, passwordHash PasswordHash) (*User, error) {
	if age < 18 {
		return nil, ErrUserTooYoung
	}

	return &User{
		ID:        NewUserID(),
		FullName:  fullName,
		Age:       age,
		IsMarried: isMarried,
		Password:  passwordHash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
