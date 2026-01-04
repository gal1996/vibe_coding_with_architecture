package persistence

import (
	"context"
	"errors"
	"sync"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
)

// MemoryWarehouseRepository is an in-memory implementation of WarehouseRepository
type MemoryWarehouseRepository struct {
	mu         sync.RWMutex
	warehouses map[string]*entity.Warehouse
}

// NewMemoryWarehouseRepository creates a new memory warehouse repository
func NewMemoryWarehouseRepository() repository.WarehouseRepository {
	return &MemoryWarehouseRepository{
		warehouses: make(map[string]*entity.Warehouse),
	}
}

func (r *MemoryWarehouseRepository) Create(ctx context.Context, warehouse *entity.Warehouse) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.warehouses[warehouse.ID]; exists {
		return errors.New("warehouse already exists")
	}

	r.warehouses[warehouse.ID] = warehouse
	return nil
}

func (r *MemoryWarehouseRepository) Update(ctx context.Context, warehouse *entity.Warehouse) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.warehouses[warehouse.ID]; !exists {
		return errors.New("warehouse not found")
	}

	r.warehouses[warehouse.ID] = warehouse
	return nil
}

func (r *MemoryWarehouseRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.warehouses[id]; !exists {
		return errors.New("warehouse not found")
	}

	delete(r.warehouses, id)
	return nil
}

func (r *MemoryWarehouseRepository) FindByID(ctx context.Context, id string) (*entity.Warehouse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	warehouse, exists := r.warehouses[id]
	if !exists {
		return nil, errors.New("warehouse not found")
	}

	// Return a copy to prevent external modifications
	copy := *warehouse
	return &copy, nil
}

func (r *MemoryWarehouseRepository) FindAll(ctx context.Context) ([]*entity.Warehouse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Warehouse
	for _, warehouse := range r.warehouses {
		copy := *warehouse
		result = append(result, &copy)
	}

	return result, nil
}