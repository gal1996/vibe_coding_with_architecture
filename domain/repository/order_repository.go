package repository

import (
	"context"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
)

// OrderRepository defines the interface for order persistence
type OrderRepository interface {
	// Create creates a new order
	Create(ctx context.Context, order *entity.Order) error

	// FindByID finds an order by its ID
	FindByID(ctx context.Context, id string) (*entity.Order, error)

	// FindByUserID finds orders by user ID
	FindByUserID(ctx context.Context, userID string) ([]*entity.Order, error)

	// Update updates an order
	Update(ctx context.Context, order *entity.Order) error
}