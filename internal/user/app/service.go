package app

import (
	"context"
	"github.com/BlackRRR/Irtea-test/internal/user/domain"
)

type UserService struct {
	userRepo  UserRepo
	hasher    PasswordHasher
	txManager TxManager
}

func NewUserService(userRepo UserRepo, hasher PasswordHasher, txManager TxManager) *UserService {
	return &UserService{
		userRepo:  userRepo,
		hasher:    hasher,
		txManager: txManager,
	}
}

func (s *UserService) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	fullName, err := domain.NewFullName(input.FirstName, input.LastName)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	passwordHash, err := domain.NewPasswordHash(hashedPassword)
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(fullName, input.Age, input.IsMarried, passwordHash)
	if err != nil {
		return nil, err
	}

	existingUser, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	err = s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		return s.userRepo.Create(txCtx, user)
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *UserService) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if err := s.hasher.Verify(user.Password.Value(), password); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	return user, nil
}
