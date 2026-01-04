package entity

import (
	"errors"
	"time"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusDelivered OrderStatus = "delivered"
)

// OrderItem represents a single item in an order
type OrderItem struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
	Subtotal    int    `json:"subtotal"`
}

// Order represents an order in the system
type Order struct {
	ID          string        `json:"id"`
	UserID      string        `json:"user_id"`
	Items       []OrderItem   `json:"items"`
	TotalPrice  int           `json:"total_price"`
	ShippingFee int           `json:"shipping_fee"`
	Status      OrderStatus   `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// NewOrder creates a new order entity
func NewOrder(id, userID string) (*Order, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	now := time.Now()
	return &Order{
		ID:          id,
		UserID:      userID,
		Items:       []OrderItem{},
		TotalPrice:  0,
		ShippingFee: 0,
		Status:      OrderStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// AddItem adds an item to the order
func (o *Order) AddItem(productID, productName string, quantity, price int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if price < 0 {
		return errors.New("price cannot be negative")
	}

	item := OrderItem{
		ProductID:   productID,
		ProductName: productName,
		Quantity:    quantity,
		Price:       price,
		Subtotal:    price * quantity,
	}

	o.Items = append(o.Items, item)
	o.calculateTotal()
	o.UpdatedAt = time.Now()
	return nil
}

// calculateTotal recalculates the total price of the order with tax and shipping
func (o *Order) calculateTotal() {
	// Calculate subtotal (before tax)
	subtotal := 0
	for _, item := range o.Items {
		subtotal += item.Subtotal
	}

	// Apply 10% consumption tax
	taxAmount := subtotal / 10 // 10% tax
	subtotalWithTax := subtotal + taxAmount

	// Calculate shipping fee
	// Free shipping for orders >= 5000 yen (before tax)
	if subtotal >= 5000 {
		o.ShippingFee = 0
	} else {
		o.ShippingFee = 500
	}

	// Calculate final total (tax included + shipping)
	o.TotalPrice = subtotalWithTax + o.ShippingFee
}

// Confirm confirms the order
func (o *Order) Confirm() error {
	if o.Status != OrderStatusPending {
		return errors.New("only pending orders can be confirmed")
	}
	if len(o.Items) == 0 {
		return errors.New("cannot confirm an empty order")
	}
	o.Status = OrderStatusConfirmed
	o.UpdatedAt = time.Now()
	return nil
}

// Cancel cancels the order
func (o *Order) Cancel() error {
	if o.Status == OrderStatusDelivered {
		return errors.New("cannot cancel delivered orders")
	}
	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
	return nil
}

// GetSubtotal returns the subtotal before tax and shipping
func (o *Order) GetSubtotal() int {
	subtotal := 0
	for _, item := range o.Items {
		subtotal += item.Subtotal
	}
	return subtotal
}

// GetTaxAmount returns the tax amount (10% of subtotal)
func (o *Order) GetTaxAmount() int {
	return o.GetSubtotal() / 10
}

// CalculateShippingFee determines the shipping fee based on subtotal
func (o *Order) CalculateShippingFee() int {
	if o.GetSubtotal() >= 5000 {
		return 0
	}
	return 500
}