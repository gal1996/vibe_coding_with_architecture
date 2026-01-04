package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
)

// OrderService handles domain logic related to orders
type OrderService struct {
	productRepo repository.ProductRepository
	orderRepo   repository.OrderRepository
}

// NewOrderService creates a new order service
func NewOrderService(productRepo repository.ProductRepository, orderRepo repository.OrderRepository) *OrderService {
	return &OrderService{
		productRepo: productRepo,
		orderRepo:   orderRepo,
	}
}

// OrderRequest represents a request to create an order
type OrderRequest struct {
	ProductID string
	Quantity  int
}

// ProcessOrder processes an order atomically
func (s *OrderService) ProcessOrder(ctx context.Context, userID string, requests []OrderRequest) (*entity.Order, error) {
	// Start a transaction
	tx, err := s.productRepo.BeginTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	productRepo := tx.GetProductRepository()

	// Create new order
	orderID := generateOrderID() // This would be implemented with a proper ID generator
	order, err := entity.NewOrder(orderID, userID)
	if err != nil {
		return nil, err
	}

	// Process each item in the order
	for _, req := range requests {
		// Check product exists and has sufficient stock
		product, err := productRepo.FindByID(ctx, req.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product not found: %s", req.ProductID)
		}

		// Check stock availability
		if !product.CanFulfillOrder(req.Quantity) {
			return nil, fmt.Errorf("insufficient stock for product %s: requested %d, available %d",
				product.Name, req.Quantity, product.Stock)
		}

		// Reduce stock atomically
		err = product.ReduceStock(req.Quantity)
		if err != nil {
			return nil, err
		}

		// Update stock in database
		err = productRepo.Update(ctx, product)
		if err != nil {
			return nil, fmt.Errorf("failed to update product stock: %w", err)
		}

		// Add item to order
		err = order.AddItem(product.ID, product.Name, req.Quantity, product.Price)
		if err != nil {
			return nil, err
		}
	}

	// Confirm the order
	err = order.Confirm()
	if err != nil {
		return nil, err
	}

	// Save the order
	err = s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}

// ValidateOrderItems validates that all requested items can be fulfilled
func (s *OrderService) ValidateOrderItems(ctx context.Context, requests []OrderRequest) error {
	for _, req := range requests {
		if req.Quantity <= 0 {
			return errors.New("quantity must be positive")
		}

		product, err := s.productRepo.FindByID(ctx, req.ProductID)
		if err != nil {
			return fmt.Errorf("product not found: %s", req.ProductID)
		}

		if !product.CanFulfillOrder(req.Quantity) {
			return fmt.Errorf("insufficient stock for product %s", product.Name)
		}
	}
	return nil
}

// generateOrderID generates a unique order ID
// This is a placeholder implementation
func generateOrderID() string {
	// In a real implementation, this would use a proper ID generator
	// For now, we'll use a simple timestamp-based ID
	return fmt.Sprintf("ORD-%d", time.Now().Unix())
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