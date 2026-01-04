package persistence

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
)

// MemoryStockRepository is an in-memory implementation of StockRepository
type MemoryStockRepository struct {
	mu     sync.RWMutex
	stocks map[string]*entity.Stock
}

// NewMemoryStockRepository creates a new memory stock repository
func NewMemoryStockRepository() repository.StockRepository {
	return &MemoryStockRepository{
		stocks: make(map[string]*entity.Stock),
	}
}

func (r *MemoryStockRepository) Create(ctx context.Context, stock *entity.Stock) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.stocks[stock.ID]; exists {
		return errors.New("stock already exists")
	}

	r.stocks[stock.ID] = stock
	return nil
}

func (r *MemoryStockRepository) Update(ctx context.Context, stock *entity.Stock) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.stocks[stock.ID]; !exists {
		return errors.New("stock not found")
	}

	r.stocks[stock.ID] = stock
	return nil
}

func (r *MemoryStockRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.stocks[id]; !exists {
		return errors.New("stock not found")
	}

	delete(r.stocks, id)
	return nil
}

func (r *MemoryStockRepository) FindByID(ctx context.Context, id string) (*entity.Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stock, exists := r.stocks[id]
	if !exists {
		return nil, errors.New("stock not found")
	}

	// Return a copy to prevent external modifications
	copy := *stock
	return &copy, nil
}

func (r *MemoryStockRepository) FindByProductID(ctx context.Context, productID string) ([]*entity.Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Stock
	for _, stock := range r.stocks {
		if stock.ProductID == productID {
			copy := *stock
			result = append(result, &copy)
		}
	}

	return result, nil
}

func (r *MemoryStockRepository) FindByWarehouseID(ctx context.Context, warehouseID string) ([]*entity.Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Stock
	for _, stock := range r.stocks {
		if stock.WarehouseID == warehouseID {
			copy := *stock
			result = append(result, &copy)
		}
	}

	return result, nil
}

func (r *MemoryStockRepository) FindByProductAndWarehouse(ctx context.Context, productID, warehouseID string) (*entity.Stock, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, stock := range r.stocks {
		if stock.ProductID == productID && stock.WarehouseID == warehouseID {
			copy := *stock
			return &copy, nil
		}
	}

	return nil, fmt.Errorf("stock not found for product %s in warehouse %s", productID, warehouseID)
}

// BeginTransaction begins a new transaction
func (r *MemoryStockRepository) BeginTransaction(ctx context.Context) (repository.StockTransaction, error) {
	return &MemoryStockTransaction{
		repo:           r,
		originalStocks: r.cloneStocks(),
	}, nil
}

// cloneStocks creates a deep copy of all stocks
func (r *MemoryStockRepository) cloneStocks() map[string]*entity.Stock {
	r.mu.RLock()
	defer r.mu.RUnlock()

	clone := make(map[string]*entity.Stock)
	for k, v := range r.stocks {
		copy := *v
		clone[k] = &copy
	}
	return clone
}

// MemoryStockTransaction represents a memory-based transaction
type MemoryStockTransaction struct {
	repo           *MemoryStockRepository
	originalStocks map[string]*entity.Stock
	committed      bool
}

func (t *MemoryStockTransaction) GetStockRepository() repository.StockRepository {
	return t.repo
}

func (t *MemoryStockTransaction) Commit() error {
	if t.committed {
		return errors.New("transaction already committed")
	}
	t.committed = true
	// Changes are already applied to the main repository
	return nil
}

func (t *MemoryStockTransaction) Rollback() error {
	if t.committed {
		return errors.New("cannot rollback committed transaction")
	}

	// Restore original stocks
	t.repo.mu.Lock()
	defer t.repo.mu.Unlock()
	t.repo.stocks = t.originalStocks
	return nil
}