package repository

import (
	"context"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
)

// UserRepository defines the interface for user persistence
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error

	// FindByID finds a user by its ID
	FindByID(ctx context.Context, id string) (*entity.User, error)

	// FindByUsername finds a user by username
	FindByUsername(ctx context.Context, username string) (*entity.User, error)

	// Update updates a user
	Update(ctx context.Context, user *entity.User) error
}