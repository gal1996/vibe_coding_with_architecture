package interactor

import (
	"context"
	"fmt"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/service"
	"github.com/gal1996/vibe_coding_with_architecture/usecase/port"
)

// WishlistUseCase handles wishlist-related business logic
type WishlistUseCase struct {
	wishlistService *service.WishlistService
	authService     port.AuthService
}

// NewWishlistUseCase creates a new wishlist use case
func NewWishlistUseCase(
	wishlistService *service.WishlistService,
	authService port.AuthService,
) *WishlistUseCase {
	return &WishlistUseCase{
		wishlistService: wishlistService,
		authService:     authService,
	}
}

// AddToWishlist adds a product to the current user's wishlist
func (uc *WishlistUseCase) AddToWishlist(ctx context.Context, productID string) error {
	// Get current user
	currentUser, err := uc.authService.GetCurrentUser(ctx)
	if err != nil {
		return fmt.Errorf("authentication required: %w", err)
	}

	// Add to wishlist
	err = uc.wishlistService.AddToWishlist(ctx, currentUser.ID, productID)
	if err != nil {
		return fmt.Errorf("failed to add to wishlist: %w", err)
	}

	return nil
}

// RemoveFromWishlist removes a product from the current user's wishlist
func (uc *WishlistUseCase) RemoveFromWishlist(ctx context.Context, productID string) error {
	// Get current user
	currentUser, err := uc.authService.GetCurrentUser(ctx)
	if err != nil {
		return fmt.Errorf("authentication required: %w", err)
	}

	// Remove from wishlist
	err = uc.wishlistService.RemoveFromWishlist(ctx, currentUser.ID, productID)
	if err != nil {
		return fmt.Errorf("failed to remove from wishlist: %w", err)
	}

	return nil
}

// GetMyWishlist gets the current user's wishlist
func (uc *WishlistUseCase) GetMyWishlist(ctx context.Context) ([]*entity.Wishlist, error) {
	// Get current user
	currentUser, err := uc.authService.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}

	// Get wishlist
	wishlist, err := uc.wishlistService.GetUserWishlist(ctx, currentUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlist: %w", err)
	}

	return wishlist, nil
}

// RecommendationResponse represents the response for recommendations
type RecommendationResponse struct {
	Products []RecommendationItem `json:"recommendations"`
}

// RecommendationItem represents a single recommendation
type RecommendationItem struct {
	ProductID    string  `json:"product_id"`
	Name         string  `json:"name"`
	Price        int     `json:"price"`
	Category     string  `json:"category"`
	Reason       string  `json:"reason"`
	Score        float64 `json:"score,omitempty"`
	TotalStock   int     `json:"total_stock"`
}

// GetRecommendations gets personalized product recommendations for the current user
func (uc *WishlistUseCase) GetRecommendations(ctx context.Context) (*RecommendationResponse, error) {
	// Get current user
	currentUser, err := uc.authService.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}

	// Get recommendations from service
	recommendations, err := uc.wishlistService.GetRecommendations(ctx, currentUser.ID, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to get recommendations: %w", err)
	}

	// Convert to response format
	response := &RecommendationResponse{
		Products: make([]RecommendationItem, len(recommendations)),
	}

	for i, rec := range recommendations {
		response.Products[i] = RecommendationItem{
			ProductID:  rec.Product.ID,
			Name:       rec.Product.Name,
			Price:      rec.Product.Price,
			Category:   rec.Product.Category,
			Reason:     rec.Reason,
			Score:      rec.Score,
			TotalStock: rec.Product.TotalStock,
		}
	}

	return response, nil
}

// CheckIsFavorite checks if a product is in the user's wishlist
func (uc *WishlistUseCase) CheckIsFavorite(ctx context.Context, productID string) (bool, error) {
	// Try to get current user (may be unauthenticated)
	currentUser, err := uc.authService.GetCurrentUser(ctx)
	if err != nil {
		// Not authenticated, return false
		return false, nil
	}

	// Check if product is in wishlist
	return uc.wishlistService.IsInWishlist(ctx, currentUser.ID, productID)
}