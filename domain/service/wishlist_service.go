package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
)

// WishlistService handles wishlist-related business logic
type WishlistService struct {
	wishlistRepo repository.WishlistRepository
	productRepo  repository.ProductRepository
	userRepo     repository.UserRepository
}

// NewWishlistService creates a new wishlist service
func NewWishlistService(
	wishlistRepo repository.WishlistRepository,
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
) *WishlistService {
	return &WishlistService{
		wishlistRepo: wishlistRepo,
		productRepo:  productRepo,
		userRepo:     userRepo,
	}
}

// AddToWishlist adds a product to user's wishlist
func (s *WishlistService) AddToWishlist(ctx context.Context, userID, productID string) error {
	// Verify user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify product exists
	_, err = s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	// Check if already in wishlist
	_, err = s.wishlistRepo.FindByUserAndProduct(ctx, userID, productID)
	if err == nil {
		return errors.New("product already in wishlist")
	}

	// Create wishlist entry
	wishlistID := fmt.Sprintf("WL-%d-%d", time.Now().Unix(), time.Now().Nanosecond())
	wishlist, err := entity.NewWishlist(wishlistID, userID, productID)
	if err != nil {
		return fmt.Errorf("failed to create wishlist entry: %w", err)
	}

	// Save to repository
	if err := s.wishlistRepo.Create(ctx, wishlist); err != nil {
		return fmt.Errorf("failed to save wishlist entry: %w", err)
	}

	return nil
}

// RemoveFromWishlist removes a product from user's wishlist
func (s *WishlistService) RemoveFromWishlist(ctx context.Context, userID, productID string) error {
	// Verify the wishlist entry exists
	_, err := s.wishlistRepo.FindByUserAndProduct(ctx, userID, productID)
	if err != nil {
		return fmt.Errorf("wishlist entry not found: %w", err)
	}

	// Delete from repository
	if err := s.wishlistRepo.Delete(ctx, userID, productID); err != nil {
		return fmt.Errorf("failed to remove from wishlist: %w", err)
	}

	return nil
}

// IsInWishlist checks if a product is in user's wishlist
func (s *WishlistService) IsInWishlist(ctx context.Context, userID, productID string) (bool, error) {
	if userID == "" {
		return false, nil // Not logged in users always get false
	}

	_, err := s.wishlistRepo.FindByUserAndProduct(ctx, userID, productID)
	if err != nil {
		// Not found is not an error in this context
		if err.Error() == "wishlist entry not found" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GetUserWishlist gets all wishlist items for a user
func (s *WishlistService) GetUserWishlist(ctx context.Context, userID string) ([]*entity.Wishlist, error) {
	return s.wishlistRepo.FindByUser(ctx, userID)
}

// RecommendationItem represents a recommended product
type RecommendationItem struct {
	Product  *entity.Product `json:"product"`
	Reason   string          `json:"reason"`
	Score    float64         `json:"score"`
}

// GetRecommendations gets product recommendations based on user's wishlist
func (s *WishlistService) GetRecommendations(ctx context.Context, userID string, limit int) ([]*RecommendationItem, error) {
	// Get user's wishlist
	wishlist, err := s.wishlistRepo.FindByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user wishlist: %w", err)
	}

	if len(wishlist) == 0 {
		// No wishlist items, return empty recommendations
		return []*RecommendationItem{}, nil
	}

	// Collect categories from wishlist products
	categoryCount := make(map[string]int)
	wishlistProductIDs := make(map[string]bool)

	for _, w := range wishlist {
		wishlistProductIDs[w.ProductID] = true
		product, err := s.productRepo.FindByID(ctx, w.ProductID)
		if err != nil {
			continue // Skip if product not found
		}
		categoryCount[product.Category]++
	}

	// Get all products and score them
	allProducts, err := s.productRepo.FindAll(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	recommendations := make([]*RecommendationItem, 0)
	for _, product := range allProducts {
		// Skip if already in wishlist
		if wishlistProductIDs[product.ID] {
			continue
		}

		// Calculate score based on category match
		score := 0.0
		if count, exists := categoryCount[product.Category]; exists {
			score = float64(count) // Higher score for categories with more wishlist items
			recommendations = append(recommendations, &RecommendationItem{
				Product: product,
				Reason:  fmt.Sprintf("同じカテゴリ「%s」の商品をお気に入りに登録されています", product.Category),
				Score:   score,
			})
		}
	}

	// Sort by score (descending) and limit
	for i := 0; i < len(recommendations)-1; i++ {
		for j := i + 1; j < len(recommendations); j++ {
			if recommendations[j].Score > recommendations[i].Score {
				recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
			}
		}
	}

	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}