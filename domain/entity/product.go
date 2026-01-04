package entity

import (
	"errors"
	"time"
)

// Product represents a product in the system
type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	Stock     int       `json:"stock"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewProduct creates a new product entity
func NewProduct(id, name string, price, stock int, category string) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if price < 0 {
		return nil, errors.New("product price cannot be negative")
	}
	if stock < 0 {
		return nil, errors.New("product stock cannot be negative")
	}
	if category == "" {
		return nil, errors.New("product category cannot be empty")
	}

	now := time.Now()
	return &Product{
		ID:        id,
		Name:      name,
		Price:     price,
		Stock:     stock,
		Category:  category,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// CanFulfillOrder checks if the product has enough stock for the requested quantity
func (p *Product) CanFulfillOrder(quantity int) bool {
	return p.Stock >= quantity
}

// ReduceStock reduces the stock by the given quantity
func (p *Product) ReduceStock(quantity int) error {
	if !p.CanFulfillOrder(quantity) {
		return errors.New("insufficient stock")
	}
	p.Stock -= quantity
	p.UpdatedAt = time.Now()
	return nil
}