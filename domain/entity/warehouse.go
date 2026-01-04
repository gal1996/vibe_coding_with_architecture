package entity

import (
	"errors"
	"fmt"
	"time"
)

// Warehouse represents a warehouse/location where products are stored
type Warehouse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewWarehouse creates a new warehouse entity
func NewWarehouse(id, name, location string) (*Warehouse, error) {
	if id == "" {
		return nil, errors.New("warehouse ID cannot be empty")
	}
	if name == "" {
		return nil, errors.New("warehouse name cannot be empty")
	}

	now := time.Now()
	return &Warehouse{
		ID:        id,
		Name:      name,
		Location:  location,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Stock represents the inventory of a product in a specific warehouse
type Stock struct {
	ID          string    `json:"id"`
	ProductID   string    `json:"product_id"`
	WarehouseID string    `json:"warehouse_id"`
	Quantity    int       `json:"quantity"`
	Reserved    int       `json:"reserved,omitempty"` // For future use: reserved quantity
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewStock creates a new stock entry
func NewStock(id, productID, warehouseID string, quantity int) (*Stock, error) {
	if id == "" {
		return nil, errors.New("stock ID cannot be empty")
	}
	if productID == "" {
		return nil, errors.New("product ID cannot be empty")
	}
	if warehouseID == "" {
		return nil, errors.New("warehouse ID cannot be empty")
	}
	if quantity < 0 {
		return nil, errors.New("quantity cannot be negative")
	}

	now := time.Now()
	return &Stock{
		ID:          id,
		ProductID:   productID,
		WarehouseID: warehouseID,
		Quantity:    quantity,
		Reserved:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// CanFulfill checks if the stock can fulfill the requested quantity
func (s *Stock) CanFulfill(requestedQuantity int) bool {
	return s.Quantity >= requestedQuantity
}

// Reduce reduces the stock quantity
func (s *Stock) Reduce(quantity int) error {
	if quantity <= 0 {
		return errors.New("reduction quantity must be positive")
	}
	if s.Quantity < quantity {
		return fmt.Errorf("insufficient stock: available=%d, requested=%d", s.Quantity, quantity)
	}
	s.Quantity -= quantity
	s.UpdatedAt = time.Now()
	return nil
}

// Add increases the stock quantity
func (s *Stock) Add(quantity int) error {
	if quantity <= 0 {
		return errors.New("addition quantity must be positive")
	}
	s.Quantity += quantity
	s.UpdatedAt = time.Now()
	return nil
}

// StockInfo represents stock information with warehouse details
type StockInfo struct {
	WarehouseID   string `json:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"`
	Quantity      int    `json:"quantity"`
}