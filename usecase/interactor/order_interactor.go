package interactor

import (
	"context"
	"fmt"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
	"github.com/gal1996/vibe_coding_with_architecture/domain/service"
	"github.com/gal1996/vibe_coding_with_architecture/usecase/port"
)

// OrderUseCase implements the order use cases
type OrderUseCase struct {
	orderRepo      repository.OrderRepository
	productRepo    repository.ProductRepository
	orderService   *service.OrderService
	authService    port.AuthService
	paymentService port.PaymentService
}

// NewOrderUseCase creates a new order use case
func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	orderService *service.OrderService,
	authService port.AuthService,
	paymentService port.PaymentService,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:      orderRepo,
		productRepo:    productRepo,
		orderService:   orderService,
		authService:    authService,
		paymentService: paymentService,
	}
}

// CreateOrderInput represents the input for creating an order
type CreateOrderInput struct {
	Items []OrderItemInput
}

// OrderItemInput represents an item in an order input
type OrderItemInput struct {
	ProductID string
	Quantity  int
}

// CreateOrder creates a new order
func (uc *OrderUseCase) CreateOrder(ctx context.Context, input CreateOrderInput) (*entity.Order, error) {
	// Get current user
	currentUser, err := uc.authService.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}

	// Convert input to domain service request
	requests := make([]service.OrderRequest, len(input.Items))
	for i, item := range input.Items {
		requests[i] = service.OrderRequest{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	// Create pending order with stock validation (but without reducing stock)
	order, err := uc.orderService.ProcessOrder(ctx, currentUser.ID, requests)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Process payment before confirming the order and reducing stock
	paymentSuccess, err := uc.paymentService.ProcessPayment(ctx, order.TotalPrice, currentUser.ID, order.ID)
	if err != nil {
		// If payment processing fails (system error), mark order as payment failed
		order.FailPayment()
		uc.orderRepo.Update(ctx, order)
		return nil, fmt.Errorf("payment processing error: %w", err)
	}

	if !paymentSuccess {
		// If payment is declined, mark order as payment failed
		order.FailPayment()
		err = uc.orderRepo.Update(ctx, order)
		if err != nil {
			return nil, fmt.Errorf("failed to update order status: %w", err)
		}
		return nil, fmt.Errorf("payment declined for order %s", order.ID)
	}

	// Payment successful, now confirm the order and reduce stock atomically
	err = uc.orderService.ConfirmOrderAndReduceStock(ctx, order)
	if err != nil {
		// If stock reduction fails, mark order as payment failed
		// (Though payment succeeded, we cannot fulfill the order)
		order.FailPayment()
		uc.orderRepo.Update(ctx, order)
		return nil, fmt.Errorf("failed to confirm order after payment: %w", err)
	}

	// Complete the order
	err = order.Complete()
	if err != nil {
		return nil, fmt.Errorf("failed to complete order: %w", err)
	}

	// Update order status in repository
	err = uc.orderRepo.Update(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return order, nil
}

// GetOrder retrieves an order by ID
func (uc *OrderUseCase) GetOrder(ctx context.Context, orderID string) (*entity.Order, error) {
	// Get current user
	currentUser, err := uc.authService.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}

	// Get order
	order, err := uc.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Check if user owns the order or is admin
	if order.UserID != currentUser.ID && !currentUser.IsAdmin {
		return nil, fmt.Errorf("permission denied: cannot access order")
	}

	return order, nil
}

// ListUserOrders lists orders for the current user
func (uc *OrderUseCase) ListUserOrders(ctx context.Context) ([]*entity.Order, error) {
	// Get current user
	currentUser, err := uc.authService.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}

	// Get user's orders
	orders, err := uc.orderRepo.FindByUserID(ctx, currentUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	return orders, nil
}
