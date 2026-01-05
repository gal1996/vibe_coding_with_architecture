package repository

import (
	"context"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
)

// CouponRepository defines the interface for coupon persistence
type CouponRepository interface {
	// Create creates a new coupon
	Create(ctx context.Context, coupon *entity.Coupon) error

	// FindByCode finds a coupon by its code
	FindByCode(ctx context.Context, code string) (*entity.Coupon, error)

	// FindByID finds a coupon by its ID
	FindByID(ctx context.Context, id string) (*entity.Coupon, error)

	// FindAll returns all coupons
	FindAll(ctx context.Context) ([]*entity.Coupon, error)

	// Update updates a coupon
	Update(ctx context.Context, coupon *entity.Coupon) error

	// Delete deletes a coupon
	Delete(ctx context.Context, id string) error
}