package service

import (
	"time"

	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/domain"
)

type TokenService struct {
	accessSecret  string
	refreshSecret string
}

func NewTokenService(accessSecret string, refreshSecret string) *TokenService {
	return &TokenService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

func (s *TokenService) IssueTokenPair(user *domain.User) (*domain.TokenPair, error) {
	return &domain.TokenPair{
		AccessToken:  "demo-access-token",
		RefreshToken: "demo-refresh-token",
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}, nil
}
