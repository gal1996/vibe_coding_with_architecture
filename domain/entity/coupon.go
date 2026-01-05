package entity

import (
	"errors"
	"time"
)

// CouponType represents the type of discount
type CouponType string

const (
	CouponTypeFixed      CouponType = "fixed"      // Fixed amount discount
	CouponTypePercentage CouponType = "percentage" // Percentage discount
)

// Coupon represents a discount coupon
type Coupon struct {
	ID           string     `json:"id"`
	Code         string     `json:"code"`         // Unique coupon code
	Description  string     `json:"description"`  // Human-readable description
	Type         CouponType `json:"type"`         // fixed or percentage
	Value        int        `json:"value"`        // Amount (yen) for fixed, percentage (0-100) for percentage
	IsActive     bool       `json:"is_active"`    // Whether the coupon is currently active
	ValidFrom    time.Time  `json:"valid_from"`   // Start of validity period
	ValidUntil   time.Time  `json:"valid_until"`  // End of validity period
	UsageLimit   int        `json:"usage_limit"`  // Maximum number of times the coupon can be used (0 = unlimited)
	UsageCount   int        `json:"usage_count"`  // Current number of times used
	MinimumOrder int        `json:"minimum_order"`// Minimum order amount required to use the coupon
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// NewCoupon creates a new coupon
func NewCoupon(id, code, description string, couponType CouponType, value int) (*Coupon, error) {
	if id == "" {
		return nil, errors.New("coupon id is required")
	}
	if code == "" {
		return nil, errors.New("coupon code is required")
	}
	if couponType != CouponTypeFixed && couponType != CouponTypePercentage {
		return nil, errors.New("invalid coupon type")
	}
	if value < 0 {
		return nil, errors.New("coupon value must be non-negative")
	}
	if couponType == CouponTypePercentage && value > 100 {
		return nil, errors.New("percentage discount cannot exceed 100")
	}

	now := time.Now()
	return &Coupon{
		ID:           id,
		Code:         code,
		Description:  description,
		Type:         couponType,
		Value:        value,
		IsActive:     true,
		ValidFrom:    now,
		ValidUntil:   now.AddDate(1, 0, 0), // Default: valid for 1 year
		UsageLimit:   0,                     // Default: unlimited usage
		UsageCount:   0,
		MinimumOrder: 0, // Default: no minimum order requirement
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// IsValid checks if the coupon is valid at the given time
func (c *Coupon) IsValid(now time.Time) bool {
	if !c.IsActive {
		return false
	}
	if now.Before(c.ValidFrom) || now.After(c.ValidUntil) {
		return false
	}
	if c.UsageLimit > 0 && c.UsageCount >= c.UsageLimit {
		return false
	}
	return true
}

// CanApplyToOrder checks if the coupon can be applied to an order with the given amount
func (c *Coupon) CanApplyToOrder(orderAmount int) bool {
	return orderAmount >= c.MinimumOrder
}

// CalculateDiscount calculates the discount amount for the given base amount
// baseAmount should be (product subtotal + tax) without shipping
func (c *Coupon) CalculateDiscount(baseAmount int) int {
	if baseAmount <= 0 {
		return 0
	}

	var discount int
	switch c.Type {
	case CouponTypeFixed:
		discount = c.Value
		if discount > baseAmount {
			discount = baseAmount // Cannot discount more than the base amount
		}
	case CouponTypePercentage:
		discount = (baseAmount * c.Value) / 100
	default:
		discount = 0
	}

	return discount
}

// IncrementUsage increments the usage count
func (c *Coupon) IncrementUsage() {
	c.UsageCount++
	c.UpdatedAt = time.Now()
}

// Deactivate deactivates the coupon
func (c *Coupon) Deactivate() {
	c.IsActive = false
	c.UpdatedAt = time.Now()
}