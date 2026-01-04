package repository

import (
	"context"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
)

// ProductRepository defines the interface for product persistence
type ProductRepository interface {
	// Create creates a new product
	Create(ctx context.Context, product *entity.Product) error

	// FindByID finds a product by its ID
	FindByID(ctx context.Context, id string) (*entity.Product, error)

	// FindAll finds all products with optional category filter
	FindAll(ctx context.Context, category string) ([]*entity.Product, error)

	// Update updates a product
	Update(ctx context.Context, product *entity.Product) error

	// UpdateStock updates the stock of a product atomically
	UpdateStock(ctx context.Context, id string, quantity int) error

	// BeginTransaction starts a new transaction
	BeginTransaction(ctx context.Context) (Transaction, error)
}

// Transaction defines the interface for database transactions
type Transaction interface {
	Commit() error
	Rollback() error
	GetProductRepository() ProductRepository
}