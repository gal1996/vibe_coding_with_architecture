package payment

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gal1996/vibe_coding_with_architecture/usecase/port"
)

// SimulatedPaymentService simulates an external payment gateway
type SimulatedPaymentService struct {
	successRate float64 // Success rate (0.0 to 1.0)
}

// NewSimulatedPaymentService creates a new simulated payment service
func NewSimulatedPaymentService() *SimulatedPaymentService {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	return &SimulatedPaymentService{
		successRate: 0.9, // 90% success rate
	}
}

// ProcessPayment simulates processing a payment
func (s *SimulatedPaymentService) ProcessPayment(ctx context.Context, amount int, userID string, orderID string) (bool, error) {
	// Log the payment attempt
	log.Printf("Processing payment: OrderID=%s, UserID=%s, Amount=%d", orderID, userID, amount)

	// Simulate network delay (50-500ms)
	delay := time.Duration(50+rand.Intn(450)) * time.Millisecond
	time.Sleep(delay)

	// Check if context is still valid
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	// Simulate payment success/failure based on success rate
	randomValue := rand.Float64()
	success := randomValue < s.successRate

	// Generate mock transaction ID for successful payments
	if success {
		transactionID := fmt.Sprintf("TXN-%d-%s", time.Now().Unix(), orderID)
		log.Printf("Payment successful: TransactionID=%s", transactionID)
	} else {
		log.Printf("Payment failed: OrderID=%s (Random value: %.2f >= Success rate: %.2f)",
			orderID, randomValue, s.successRate)
	}

	return success, nil
}

// SetSuccessRate allows changing the success rate for testing
func (s *SimulatedPaymentService) SetSuccessRate(rate float64) {
	if rate < 0 {
		rate = 0
	} else if rate > 1 {
		rate = 1
	}
	s.successRate = rate
	log.Printf("Payment service success rate set to: %.2f%%", rate*100)
}

// Ensure SimulatedPaymentService implements port.PaymentService
var _ port.PaymentService = (*SimulatedPaymentService)(nil)