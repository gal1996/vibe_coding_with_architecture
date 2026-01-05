package entity

import (
	"errors"
	"time"
)

// Wishlist represents a user's favorite product
type Wishlist struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
}

// NewWishlist creates a new wishlist entry
func NewWishlist(id, userID, productID string) (*Wishlist, error) {
	if id == "" {
		return nil, errors.New("wishlist id is required")
	}
	if userID == "" {
		return nil, errors.New("user id is required")
	}
	if productID == "" {
		return nil, errors.New("product id is required")
	}

	return &Wishlist{
		ID:        id,
		UserID:    userID,
		ProductID: productID,
		CreatedAt: time.Now(),
	}, nil
}

// Validate validates the wishlist entry
func (w *Wishlist) Validate() error {
	if w.UserID == "" {
		return errors.New("user id is required")
	}
	if w.ProductID == "" {
		return errors.New("product id is required")
	}
	return nil
}