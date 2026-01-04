package persistence

import (
	"context"
	"errors"
	"sync"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
)

// MemoryUserRepository is an in-memory implementation of UserRepository
type MemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*entity.User
}

// NewMemoryUserRepository creates a new in-memory user repository
func NewMemoryUserRepository() repository.UserRepository {
	return &MemoryUserRepository{
		users: make(map[string]*entity.User),
	}
}

// Create creates a new user
func (r *MemoryUserRepository) Create(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user ID already exists
	if _, exists := r.users[user.ID]; exists {
		return errors.New("user already exists")
	}

	// Check if username already exists
	for _, u := range r.users {
		if u.Username == user.Username {
			return errors.New("username already taken")
		}
	}

	// Create a copy to avoid external modifications
	userCopy := *user
	r.users[user.ID] = &userCopy
	return nil
}

// FindByID finds a user by its ID
func (r *MemoryUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Return a copy to avoid external modifications
	userCopy := *user
	return &userCopy, nil
}

// FindByUsername finds a user by username
func (r *MemoryUserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			// Return a copy to avoid external modifications
			userCopy := *user
			return &userCopy, nil
		}
	}

	return nil, errors.New("user not found")
}

// Update updates a user
func (r *MemoryUserRepository) Update(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	// Check if username is being changed and is already taken
	existingUser := r.users[user.ID]
	if existingUser.Username != user.Username {
		for _, u := range r.users {
			if u.ID != user.ID && u.Username == user.Username {
				return errors.New("username already taken")
			}
		}
	}

	// Create a copy to avoid external modifications
	userCopy := *user
	r.users[user.ID] = &userCopy
	return nil
}