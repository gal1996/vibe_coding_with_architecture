package persistence

import (
	"context"
	"errors"
	"sync"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
)

// MemoryCouponRepository is an in-memory implementation of CouponRepository
type MemoryCouponRepository struct {
	mu      sync.RWMutex
	coupons map[string]*entity.Coupon
}

// NewMemoryCouponRepository creates a new memory coupon repository
func NewMemoryCouponRepository() repository.CouponRepository {
	return &MemoryCouponRepository{
		coupons: make(map[string]*entity.Coupon),
	}
}

func (r *MemoryCouponRepository) Create(ctx context.Context, coupon *entity.Coupon) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.coupons[coupon.ID]; exists {
		return errors.New("coupon already exists")
	}

	// Check if code is unique
	for _, existing := range r.coupons {
		if existing.Code == coupon.Code {
			return errors.New("coupon code already exists")
		}
	}

	r.coupons[coupon.ID] = coupon
	return nil
}

func (r *MemoryCouponRepository) FindByCode(ctx context.Context, code string) (*entity.Coupon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, coupon := range r.coupons {
		if coupon.Code == code {
			// Return a copy to prevent external modifications
			copy := *coupon
			return &copy, nil
		}
	}

	return nil, errors.New("coupon not found")
}

func (r *MemoryCouponRepository) FindByID(ctx context.Context, id string) (*entity.Coupon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	coupon, exists := r.coupons[id]
	if !exists {
		return nil, errors.New("coupon not found")
	}

	// Return a copy to prevent external modifications
	copy := *coupon
	return &copy, nil
}

func (r *MemoryCouponRepository) FindAll(ctx context.Context) ([]*entity.Coupon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Coupon
	for _, coupon := range r.coupons {
		copy := *coupon
		result = append(result, &copy)
	}

	return result, nil
}

func (r *MemoryCouponRepository) Update(ctx context.Context, coupon *entity.Coupon) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.coupons[coupon.ID]; !exists {
		return errors.New("coupon not found")
	}

	r.coupons[coupon.ID] = coupon
	return nil
}

func (r *MemoryCouponRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.coupons[id]; !exists {
		return errors.New("coupon not found")
	}

	delete(r.coupons, id)
	return nil
}