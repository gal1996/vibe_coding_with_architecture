package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
)

// CouponService handles coupon-related business logic
type CouponService struct {
	couponRepo repository.CouponRepository
}

// NewCouponService creates a new coupon service
func NewCouponService(couponRepo repository.CouponRepository) *CouponService {
	return &CouponService{
		couponRepo: couponRepo,
	}
}

// ValidateAndGetCoupon validates a coupon code and returns the coupon if valid
func (s *CouponService) ValidateAndGetCoupon(ctx context.Context, code string) (*entity.Coupon, error) {
	if code == "" {
		return nil, nil // No coupon to apply
	}

	coupon, err := s.couponRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("invalid coupon code: %s", code)
	}

	// Check if coupon is valid
	if !coupon.IsValid(time.Now()) {
		return nil, fmt.Errorf("coupon %s is not valid or has expired", code)
	}

	return coupon, nil
}

// ApplyCoupon applies a coupon to an order and calculates the discount
func (s *CouponService) ApplyCoupon(ctx context.Context, coupon *entity.Coupon, baseAmount int) (int, error) {
	if coupon == nil {
		return 0, nil // No coupon to apply
	}

	// Check if order meets minimum requirement
	if !coupon.CanApplyToOrder(baseAmount) {
		return 0, fmt.Errorf("order amount does not meet minimum requirement for coupon %s (minimum: %d yen)",
			coupon.Code, coupon.MinimumOrder)
	}

	// Calculate discount
	discount := coupon.CalculateDiscount(baseAmount)

	// Increment usage count
	coupon.IncrementUsage()
	err := s.couponRepo.Update(ctx, coupon)
	if err != nil {
		return 0, fmt.Errorf("failed to update coupon usage: %w", err)
	}

	return discount, nil
}

// RollbackCouponUsage rolls back the usage count of a coupon (used when payment fails)
func (s *CouponService) RollbackCouponUsage(ctx context.Context, couponCode string) error {
	if couponCode == "" {
		return nil // No coupon was used
	}

	coupon, err := s.couponRepo.FindByCode(ctx, couponCode)
	if err != nil {
		// Log the error but don't fail the rollback
		// In production, this should be logged properly
		return nil
	}

	// Decrement usage count
	if coupon.UsageCount > 0 {
		coupon.UsageCount--
		coupon.UpdatedAt = time.Now()
		_ = s.couponRepo.Update(ctx, coupon)
	}

	return nil
}