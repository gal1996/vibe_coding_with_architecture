package entity

import (
	"errors"
	"time"
)

// Product represents a product in the system
type Product struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Price     int         `json:"price"`
	Category  string      `json:"category"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	// Stock is now managed separately through Stock entities
	Stocks    []StockInfo `json:"stocks,omitempty"`     // Stock information by warehouse
	TotalStock int        `json:"total_stock,omitempty"` // Calculated total across all warehouses
}

// NewProduct creates a new product entity
func NewProduct(id, name string, price int, category string) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if price < 0 {
		return nil, errors.New("product price cannot be negative")
	}
	if category == "" {
		return nil, errors.New("product category cannot be empty")
	}

	now := time.Now()
	return &Product{
		ID:        id,
		Name:      name,
		Price:     price,
		Category:  category,
		CreatedAt: now,
		UpdatedAt: now,
		Stocks:    []StockInfo{},
		TotalStock: 0,
	}, nil
}

// CalculateTotalStock calculates the total stock across all warehouses
func (p *Product) CalculateTotalStock() int {
	total := 0
	for _, stock := range p.Stocks {
		total += stock.Quantity
	}
	p.TotalStock = total
	return total
}

// CanFulfillOrder checks if the product has enough total stock for the requested quantity
func (p *Product) CanFulfillOrder(quantity int) bool {
	return p.TotalStock >= quantity
}

// AddStockInfo adds stock information for a warehouse
func (p *Product) AddStockInfo(info StockInfo) {
	p.Stocks = append(p.Stocks, info)
	p.CalculateTotalStock()
}