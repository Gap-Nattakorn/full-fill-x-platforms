package repository

import "context"

type TokenRepository interface {
	BlacklistRefreshToken(ctx context.Context, token string) error
	IsRefreshTokenBlacklisted(ctx context.Context, token string) (bool, error)
}
