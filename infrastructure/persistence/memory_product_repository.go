package persistence

import (
	"context"
	"errors"
	"sync"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
)

// MemoryProductRepository is an in-memory implementation of ProductRepository
type MemoryProductRepository struct {
	mu       sync.RWMutex
	products map[string]*entity.Product
}

// NewMemoryProductRepository creates a new in-memory product repository
func NewMemoryProductRepository() *MemoryProductRepository {
	return &MemoryProductRepository{
		products: make(map[string]*entity.Product),
	}
}

// Create creates a new product
func (r *MemoryProductRepository) Create(ctx context.Context, product *entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; exists {
		return errors.New("product already exists")
	}

	// Create a copy to avoid external modifications
	productCopy := *product
	r.products[product.ID] = &productCopy
	return nil
}

// FindByID finds a product by its ID
func (r *MemoryProductRepository) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, exists := r.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}

	// Return a copy to avoid external modifications
	productCopy := *product
	return &productCopy, nil
}

// FindAll finds all products with optional category filter
func (r *MemoryProductRepository) FindAll(ctx context.Context, category string) ([]*entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Product
	for _, product := range r.products {
		if category == "" || product.Category == category {
			// Create a copy to avoid external modifications
			productCopy := *product
			result = append(result, &productCopy)
		}
	}
	return result, nil
}

// Update updates a product
func (r *MemoryProductRepository) Update(ctx context.Context, product *entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.products[product.ID]; !exists {
		return errors.New("product not found")
	}

	// Create a copy to avoid external modifications
	productCopy := *product
	r.products[product.ID] = &productCopy
	return nil
}

// Note: Stock management is now handled through StockService and StockRepository
// Products themselves don't maintain stock counts anymore

// BeginTransaction starts a new transaction
func (r *MemoryProductRepository) BeginTransaction(ctx context.Context) (repository.Transaction, error) {
	return &MemoryTransaction{
		productRepo: r,
		committed:   false,
		rollback:    false,
	}, nil
}

// MemoryTransaction represents an in-memory transaction
type MemoryTransaction struct {
	productRepo *MemoryProductRepository
	committed   bool
	rollback    bool
}

// Commit commits the transaction
func (t *MemoryTransaction) Commit() error {
	if t.rollback {
		return errors.New("transaction already rolled back")
	}
	if t.committed {
		return errors.New("transaction already committed")
	}
	t.committed = true
	return nil
}

// Rollback rolls back the transaction
func (t *MemoryTransaction) Rollback() error {
	if t.committed {
		return errors.New("transaction already committed")
	}
	if t.rollback {
		return errors.New("transaction already rolled back")
	}
	t.rollback = true
	return nil
}

// GetProductRepository returns the product repository for this transaction
func (t *MemoryTransaction) GetProductRepository() repository.ProductRepository {
	return t.productRepo
}