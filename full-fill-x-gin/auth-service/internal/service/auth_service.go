package service

import (
	"context"
	"errors"

	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/domain"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/repository"
	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/security"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	users          repository.UserRepository
	tokens         *TokenService
	passwordHasher security.PasswordHasher
}

func NewAuthService(users repository.UserRepository, tokens *TokenService, passwordHasher security.PasswordHasher) *AuthService {
	return &AuthService{
		users:          users,
		tokens:         tokens,
		passwordHasher: passwordHasher,
	}
}

func (s *AuthService) Register(ctx context.Context, email string, password string) (*domain.User, *domain.TokenPair, error) {
	passwordHash, err := s.passwordHasher.Hash(password)
	if err != nil {
		return nil, nil, err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: passwordHash,
		Role:         domain.RoleCustomer,
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	tokenPair, err := s.tokens.IssueTokenPair(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokenPair, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (*domain.User, *domain.TokenPair, error) {
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	if err := s.passwordHasher.Compare(user.PasswordHash, password); err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	tokenPair, err := s.tokens.IssueTokenPair(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokenPair, nil
}
