package persistence

import (
	"context"
	"errors"
	"sync"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
)

// MemoryWishlistRepository is an in-memory implementation of WishlistRepository
type MemoryWishlistRepository struct {
	mu        sync.RWMutex
	wishlists map[string]*entity.Wishlist // key: wishlist ID
}

// NewMemoryWishlistRepository creates a new in-memory wishlist repository
func NewMemoryWishlistRepository() repository.WishlistRepository {
	return &MemoryWishlistRepository{
		wishlists: make(map[string]*entity.Wishlist),
	}
}

// Create adds a new wishlist entry
func (r *MemoryWishlistRepository) Create(ctx context.Context, wishlist *entity.Wishlist) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if wishlist ID already exists
	if _, exists := r.wishlists[wishlist.ID]; exists {
		return errors.New("wishlist entry already exists")
	}

	// Check if user-product pair already exists
	for _, w := range r.wishlists {
		if w.UserID == wishlist.UserID && w.ProductID == wishlist.ProductID {
			return errors.New("product already in user's wishlist")
		}
	}

	// Create a copy to avoid external modifications
	wishlistCopy := *wishlist
	r.wishlists[wishlist.ID] = &wishlistCopy
	return nil
}

// Delete removes a wishlist entry by user and product ID
func (r *MemoryWishlistRepository) Delete(ctx context.Context, userID, productID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Find and delete the wishlist entry
	for id, w := range r.wishlists {
		if w.UserID == userID && w.ProductID == productID {
			delete(r.wishlists, id)
			return nil
		}
	}

	return errors.New("wishlist entry not found")
}

// FindByUserAndProduct checks if a specific product is in user's wishlist
func (r *MemoryWishlistRepository) FindByUserAndProduct(ctx context.Context, userID, productID string) (*entity.Wishlist, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, w := range r.wishlists {
		if w.UserID == userID && w.ProductID == productID {
			// Return a copy to avoid external modifications
			wishlistCopy := *w
			return &wishlistCopy, nil
		}
	}

	return nil, errors.New("wishlist entry not found")
}

// FindByUser gets all wishlist entries for a user
func (r *MemoryWishlistRepository) FindByUser(ctx context.Context, userID string) ([]*entity.Wishlist, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userWishlists []*entity.Wishlist
	for _, w := range r.wishlists {
		if w.UserID == userID {
			// Create a copy to avoid external modifications
			wishlistCopy := *w
			userWishlists = append(userWishlists, &wishlistCopy)
		}
	}

	return userWishlists, nil
}

// FindByProduct gets all wishlist entries for a product
func (r *MemoryWishlistRepository) FindByProduct(ctx context.Context, productID string) ([]*entity.Wishlist, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var productWishlists []*entity.Wishlist
	for _, w := range r.wishlists {
		if w.ProductID == productID {
			// Create a copy to avoid external modifications
			wishlistCopy := *w
			productWishlists = append(productWishlists, &wishlistCopy)
		}
	}

	return productWishlists, nil
}

// CountByUser gets the count of wishlist items for a user
func (r *MemoryWishlistRepository) CountByUser(ctx context.Context, userID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, w := range r.wishlists {
		if w.UserID == userID {
			count++
		}
	}

	return count, nil
}