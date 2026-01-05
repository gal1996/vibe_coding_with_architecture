package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gal1996/vibe_coding_with_architecture/usecase/interactor"
)

// WishlistHandler handles HTTP requests for wishlist operations
type WishlistHandler struct {
	wishlistUseCase *interactor.WishlistUseCase
}

// NewWishlistHandler creates a new wishlist handler
func NewWishlistHandler(wishlistUseCase *interactor.WishlistUseCase) *WishlistHandler {
	return &WishlistHandler{
		wishlistUseCase: wishlistUseCase,
	}
}

// AddToWishlist handles POST /wishlist/:product_id
func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	productID := c.Param("product_id")

	err := h.wishlistUseCase.AddToWishlist(c.Request.Context(), productID)
	if err != nil {
		// Check if it's an authentication error
		if err.Error() == "authentication required" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}
		// Check if product already in wishlist
		if err.Error() == "product already in wishlist" {
			c.JSON(http.StatusConflict, gin.H{"error": "Product already in wishlist"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Product added to wishlist",
		"product_id": productID,
	})
}

// RemoveFromWishlist handles DELETE /wishlist/:product_id
func (h *WishlistHandler) RemoveFromWishlist(c *gin.Context) {
	productID := c.Param("product_id")

	err := h.wishlistUseCase.RemoveFromWishlist(c.Request.Context(), productID)
	if err != nil {
		// Check if it's an authentication error
		if err.Error() == "authentication required" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}
		// Check if product not in wishlist
		if err.Error() == "wishlist entry not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not in wishlist"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Product removed from wishlist",
		"product_id": productID,
	})
}

// GetMyWishlist handles GET /wishlist
func (h *WishlistHandler) GetMyWishlist(c *gin.Context) {
	wishlist, err := h.wishlistUseCase.GetMyWishlist(c.Request.Context())
	if err != nil {
		// Check if it's an authentication error
		if err.Error() == "authentication required" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wishlist": wishlist,
		"count":    len(wishlist),
	})
}

// GetRecommendations handles GET /users/me/recommendations
func (h *WishlistHandler) GetRecommendations(c *gin.Context) {
	recommendations, err := h.wishlistUseCase.GetRecommendations(c.Request.Context())
	if err != nil {
		// Check if it's an authentication error
		if err.Error() == "authentication required" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}