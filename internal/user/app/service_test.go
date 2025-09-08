package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/BlackRRR/Irtea-test/internal/user/domain"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) Delete(ctx context.Context, id domain.UserID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) Verify(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

type MockTxManager struct {
	mock.Mock
}

func (m *MockTxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	if args.Get(0) != nil {
		return args.Get(0).(func(context.Context, func(context.Context) error) error)(ctx, fn)
	}
	return fn(ctx)
}

func TestUserService_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockHasher := new(MockPasswordHasher)
	mockTx := new(MockTxManager)

	service := NewUserService(mockRepo, mockHasher, mockTx)

	input := RegisterInput{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@domain.com",
		Age:       25,
		IsMarried: false,
		Password:  "password123",
	}

	mockRepo.On("GetByEmail", mock.Anything, "john.doe@domain.com").Return(nil, domain.ErrUserNotFound)
	mockHasher.On("Hash", "password123").Return("hashedpassword", nil)
	mockTx.On("WithTx", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := service.Register(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John", user.FullName.FirstName)
	assert.Equal(t, "Doe", user.FullName.LastName)
	assert.Equal(t, 25, user.Age)
	assert.Equal(t, false, user.IsMarried)

	mockRepo.AssertExpectations(t)
	mockHasher.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestUserService_Register_TooYoung(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockHasher := new(MockPasswordHasher)
	mockTx := new(MockTxManager)

	service := NewUserService(mockRepo, mockHasher, mockTx)

	input := RegisterInput{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@domain.com",
		Age:       17,
		IsMarried: false,
		Password:  "password123",
	}

	user, err := service.Register(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, domain.ErrUserTooYoung, err)

	mockRepo.AssertNotCalled(t, "Create")
	mockHasher.AssertNotCalled(t, "Hash")
}

func TestUserService_Register_WeakPassword(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockHasher := new(MockPasswordHasher)
	mockTx := new(MockTxManager)

	service := NewUserService(mockRepo, mockHasher, mockTx)

	input := RegisterInput{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@domain.com",
		Age:       25,
		IsMarried: false,
		Password:  "123",
	}

	user, err := service.Register(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, domain.ErrInvalidPassword, err)

	mockRepo.AssertNotCalled(t, "Create")
	mockHasher.AssertNotCalled(t, "Hash")
}

func TestUserService_Register_UserAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockHasher := new(MockPasswordHasher)
	mockTx := new(MockTxManager)

	service := NewUserService(mockRepo, mockHasher, mockTx)

	existingUser := &domain.User{}
	
	input := RegisterInput{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@domain.com",
		Age:       25,
		IsMarried: false,
		Password:  "password123",
	}

	mockRepo.On("GetByEmail", mock.Anything, "john.doe@domain.com").Return(existingUser, nil)

	user, err := service.Register(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, domain.ErrUserAlreadyExists, err)

	mockRepo.AssertExpectations(t)
	mockHasher.AssertNotCalled(t, "Hash")
}