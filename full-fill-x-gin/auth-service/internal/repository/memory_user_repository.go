package repository

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/auth-service/internal/domain"
)

var ErrUserNotFound = errors.New("user not found")

type MemoryUserRepository struct {
	mu      sync.RWMutex
	byID    map[string]*domain.User
	byEmail map[string]*domain.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		byID:    make(map[string]*domain.User),
		byEmail: make(map[string]*domain.User),
	}
}

func (r *MemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	user.ID = strconv.FormatInt(now.UnixNano(), 36)
	user.CreatedAt = now
	user.UpdatedAt = now

	r.byID[user.ID] = user
	r.byEmail[user.Email] = user

	return nil
}

func (r *MemoryUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.byEmail[email]
	if !ok {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (r *MemoryUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.byID[id]
	if !ok {
		return nil, ErrUserNotFound
	}

	return user, nil
}
