package port

import (
	"context"
)

// PaymentService represents the payment gateway interface
type PaymentService interface {
	// ProcessPayment processes the payment for the given amount
	// Returns true if payment is successful, false otherwise
	ProcessPayment(ctx context.Context, amount int, userID string, orderID string) (bool, error)
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	Amount  int    `json:"amount"`
	UserID  string `json:"user_id"`
	OrderID string `json:"order_id"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	Success       bool   `json:"success"`
	TransactionID string `json:"transaction_id,omitempty"`
	ErrorMessage  string `json:"error_message,omitempty"`
}