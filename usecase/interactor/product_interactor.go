package interactor

import (
	"context"
	"errors"
	"fmt"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
	"github.com/gal1996/vibe_coding_with_architecture/usecase/port"
)

// ProductUseCase handles product-related business logic
type ProductUseCase struct {
	productRepo repository.ProductRepository
	userRepo    repository.UserRepository
	authService port.AuthService
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
	authService port.AuthService,
) *ProductUseCase {
	return &ProductUseCase{
		productRepo: productRepo,
		userRepo:    userRepo,
		authService: authService,
	}
}

// CreateProductInput represents the input for creating a product
type CreateProductInput struct {
	Name     string
	Price    int
	Stock    int
	Category string
}

// CreateProduct creates a new product (admin only)
func (uc *ProductUseCase) CreateProduct(ctx context.Context, input CreateProductInput) (*entity.Product, error) {
	// Get current user
	currentUser, err := uc.authService.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}

	// Check if user is admin
	if !currentUser.CanCreateProduct() {
		return nil, errors.New("permission denied: only admins can create products")
	}

	// Generate product ID
	productID := generateProductID()

	// Create product entity
	product, err := entity.NewProduct(productID, input.Name, input.Price, input.Stock, input.Category)
	if err != nil {
		return nil, fmt.Errorf("invalid product data: %w", err)
	}

	// Save to repository
	err = uc.productRepo.Create(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// GetProduct retrieves a product by ID
func (uc *ProductUseCase) GetProduct(ctx context.Context, productID string) (*entity.Product, error) {
	product, err := uc.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return product, nil
}

// ListProducts lists all products with optional category filter
func (uc *ProductUseCase) ListProducts(ctx context.Context, category string) ([]*entity.Product, error) {
	products, err := uc.productRepo.FindAll(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return products, nil
}

// generateProductID generates a unique product ID
func generateProductID() string {
	// In a real implementation, this would use a proper ID generator
	return fmt.Sprintf("PROD-%d", time.Now().Unix())
}

var time = struct {
	Now func() struct {
		Unix func() int64
	}
}{
	Now: func() struct {
		Unix func() int64
	} {
		return struct {
			Unix func() int64
		}{
			Unix: func() int64 {
				return 1234567890 // Placeholder
			},
		}
	},
}