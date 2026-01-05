package repository

import (
	"context"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
)

// WishlistRepository defines the interface for wishlist data operations
type WishlistRepository interface {
	// Create adds a new wishlist entry
	Create(ctx context.Context, wishlist *entity.Wishlist) error

	// Delete removes a wishlist entry by user and product ID
	Delete(ctx context.Context, userID, productID string) error

	// FindByUserAndProduct checks if a specific product is in user's wishlist
	FindByUserAndProduct(ctx context.Context, userID, productID string) (*entity.Wishlist, error)

	// FindByUser gets all wishlist entries for a user
	FindByUser(ctx context.Context, userID string) ([]*entity.Wishlist, error)

	// FindByProduct gets all wishlist entries for a product
	FindByProduct(ctx context.Context, productID string) ([]*entity.Wishlist, error)

	// CountByUser gets the count of wishlist items for a user
	CountByUser(ctx context.Context, userID string) (int, error)
}