package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
)

// OrderService handles domain logic related to orders
type OrderService struct {
	productRepo repository.ProductRepository
	orderRepo   repository.OrderRepository
	stockService *StockService
}

// NewOrderService creates a new order service
func NewOrderService(productRepo repository.ProductRepository, orderRepo repository.OrderRepository, stockService *StockService) *OrderService {
	return &OrderService{
		productRepo: productRepo,
		orderRepo:   orderRepo,
		stockService: stockService,
	}
}

// OrderRequest represents a request to create an order
type OrderRequest struct {
	ProductID string
	Quantity  int
}

// ProcessOrder creates a pending order with stock validation but without reducing stock
// Stock reduction happens after payment is confirmed
func (s *OrderService) ProcessOrder(ctx context.Context, userID string, requests []OrderRequest) (*entity.Order, error) {
	// Create new order
	orderID := generateOrderID() // This would be implemented with a proper ID generator
	order, err := entity.NewOrder(orderID, userID)
	if err != nil {
		return nil, err
	}

	// Validate and add items to order
	for _, req := range requests {
		// Check product exists
		product, err := s.productRepo.FindByID(ctx, req.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product not found: %s", req.ProductID)
		}

		// Check stock availability across all warehouses
		available, totalStock, err := s.stockService.CheckAvailability(ctx, req.ProductID, req.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to check stock availability: %w", err)
		}

		if !available {
			return nil, fmt.Errorf("insufficient stock for product %s: requested %d, available %d",
				product.Name, req.Quantity, totalStock)
		}

		// Add item to order (without reducing stock)
		err = order.AddItem(product.ID, product.Name, req.Quantity, product.Price)
		if err != nil {
			return nil, err
		}
	}

	// Keep order in pending status for payment processing
	// Don't confirm yet - wait for payment

	// Save the order
	err = s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

// ConfirmOrderAndReduceStock confirms the order and reduces stock after successful payment
func (s *OrderService) ConfirmOrderAndReduceStock(ctx context.Context, order *entity.Order) error {
	// Store allocations for potential rollback
	allocations := make(map[string][]StockAllocation)

	// Allocate stock for each item in the order
	for _, item := range order.Items {
		// Double-check stock availability
		available, totalStock, err := s.stockService.CheckAvailability(ctx, item.ProductID, item.Quantity)
		if err != nil {
			// Rollback any previous allocations
			s.rollbackAllocations(ctx, allocations)
			return fmt.Errorf("failed to check stock: %w", err)
		}

		if !available {
			// Rollback any previous allocations
			s.rollbackAllocations(ctx, allocations)
			return fmt.Errorf("insufficient stock for product %s: requested %d, available %d",
				item.ProductName, item.Quantity, totalStock)
		}

		// Allocate stock from warehouses
		itemAllocations, err := s.stockService.AllocateStock(ctx, item.ProductID, item.Quantity)
		if err != nil {
			// Rollback any previous allocations
			s.rollbackAllocations(ctx, allocations)
			return fmt.Errorf("failed to allocate stock for product %s: %w", item.ProductName, err)
		}

		allocations[item.ProductID] = itemAllocations
	}

	// Confirm the order
	err := order.Confirm()
	if err != nil {
		// Rollback all allocations if order confirmation fails
		s.rollbackAllocations(ctx, allocations)
		return err
	}

	return nil
}

// rollbackAllocations restores stock allocations in case of failure
func (s *OrderService) rollbackAllocations(ctx context.Context, allocations map[string][]StockAllocation) {
	for productID, stockAllocations := range allocations {
		_ = s.stockService.RestoreStock(ctx, productID, stockAllocations)
	}
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
	// For now, we'll use a timestamp with random component
	return fmt.Sprintf("ORD-%d-%d", time.Now().Unix(), rand.Intn(100000))
}