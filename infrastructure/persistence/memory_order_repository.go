package persistence

import (
	"context"
	"errors"
	"sync"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
)

// MemoryOrderRepository is an in-memory implementation of OrderRepository
type MemoryOrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*entity.Order
}

// NewMemoryOrderRepository creates a new in-memory order repository
func NewMemoryOrderRepository() repository.OrderRepository {
	return &MemoryOrderRepository{
		orders: make(map[string]*entity.Order),
	}
}

// Create creates a new order
func (r *MemoryOrderRepository) Create(ctx context.Context, order *entity.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.ID]; exists {
		return errors.New("order already exists")
	}

	// Create a deep copy to avoid external modifications
	orderCopy := *order
	orderCopy.Items = make([]entity.OrderItem, len(order.Items))
	copy(orderCopy.Items, order.Items)
	r.orders[order.ID] = &orderCopy
	return nil
}

// FindByID finds an order by its ID
func (r *MemoryOrderRepository) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}

	// Return a deep copy to avoid external modifications
	orderCopy := *order
	orderCopy.Items = make([]entity.OrderItem, len(order.Items))
	copy(orderCopy.Items, order.Items)
	return &orderCopy, nil
}

// FindByUserID finds orders by user ID
func (r *MemoryOrderRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Order
	for _, order := range r.orders {
		if order.UserID == userID {
			// Create a deep copy to avoid external modifications
			orderCopy := *order
			orderCopy.Items = make([]entity.OrderItem, len(order.Items))
			copy(orderCopy.Items, order.Items)
			result = append(result, &orderCopy)
		}
	}
	return result, nil
}

// Update updates an order
func (r *MemoryOrderRepository) Update(ctx context.Context, order *entity.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[order.ID]; !exists {
		return errors.New("order not found")
	}

	// Create a deep copy to avoid external modifications
	orderCopy := *order
	orderCopy.Items = make([]entity.OrderItem, len(order.Items))
	copy(orderCopy.Items, order.Items)
	r.orders[order.ID] = &orderCopy
	return nil
}