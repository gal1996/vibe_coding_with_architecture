package port

import (
	"context"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
)

// AuthService defines the interface for authentication services
type AuthService interface {
	// GenerateToken generates an authentication token for a user
	GenerateToken(user *entity.User) (string, error)

	// ValidateToken validates an authentication token and returns the user ID
	ValidateToken(token string) (string, error)

	// GetCurrentUser gets the current user from the context
	GetCurrentUser(ctx context.Context) (*entity.User, error)
}