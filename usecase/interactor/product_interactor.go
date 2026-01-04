package interactor

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
	"github.com/gal1996/vibe_coding_with_architecture/domain/service"
	"github.com/gal1996/vibe_coding_with_architecture/usecase/port"
)

// ProductUseCase handles product-related business logic
type ProductUseCase struct {
	productRepo  repository.ProductRepository
	userRepo     repository.UserRepository
	authService  port.AuthService
	stockService *service.StockService
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(
	productRepo repository.ProductRepository,
	userRepo repository.UserRepository,
	authService port.AuthService,
	stockService *service.StockService,
) *ProductUseCase {
	return &ProductUseCase{
		productRepo:  productRepo,
		userRepo:     userRepo,
		authService:  authService,
		stockService: stockService,
	}
}

// CreateProductInput represents the input for creating a product
type CreateProductInput struct {
	Name     string
	Price    int
	Category string
	// Stock is now managed through warehouse-specific allocations
	// Use StockService to add stock to specific warehouses after product creation
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

	// Create product entity (without stock - stock is managed through warehouses)
	product, err := entity.NewProduct(productID, input.Name, input.Price, input.Category)
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

// GetProduct retrieves a product by ID with stock information from all warehouses
func (uc *ProductUseCase) GetProduct(ctx context.Context, productID string) (*entity.Product, error) {
	product, err := uc.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Get stock information from all warehouses
	stockInfos, totalStock, err := uc.stockService.GetProductStockInfo(ctx, productID)
	if err != nil {
		// Log error but don't fail - product exists even if stock info unavailable
		// In production, you might want to handle this differently
		stockInfos = []entity.StockInfo{}
		totalStock = 0
	}

	// Add warehouse stock information to product
	product.Stocks = stockInfos
	product.TotalStock = totalStock

	return product, nil
}

// ListProducts lists all products with optional category filter and stock information
func (uc *ProductUseCase) ListProducts(ctx context.Context, category string) ([]*entity.Product, error) {
	products, err := uc.productRepo.FindAll(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	// Add stock information for each product
	for _, product := range products {
		stockInfos, totalStock, err := uc.stockService.GetProductStockInfo(ctx, product.ID)
		if err != nil {
			// Log error but don't fail - continue with empty stock info
			stockInfos = []entity.StockInfo{}
			totalStock = 0
		}
		product.Stocks = stockInfos
		product.TotalStock = totalStock
	}

	return products, nil
}

// generateProductID generates a unique product ID
func generateProductID() string {
	// In a real implementation, this would use a proper ID generator
	return fmt.Sprintf("PROD-%d-%d", time.Now().Unix(), rand.Intn(10000))
}