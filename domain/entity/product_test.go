package entity

import (
	"testing"
)

func TestNewProduct(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		prodName string
		price    int
		stock    int
		category string
		wantErr  bool
	}{
		{
			name:     "valid product",
			id:       "PROD-001",
			prodName: "Laptop",
			price:    1200,
			stock:    10,
			category: "Electronics",
			wantErr:  false,
		},
		{
			name:     "empty name",
			id:       "PROD-002",
			prodName: "",
			price:    100,
			stock:    5,
			category: "Test",
			wantErr:  true,
		},
		{
			name:     "negative price",
			id:       "PROD-003",
			prodName: "Test",
			price:    -100,
			stock:    5,
			category: "Test",
			wantErr:  true,
		},
		{
			name:     "negative stock",
			id:       "PROD-004",
			prodName: "Test",
			price:    100,
			stock:    -5,
			category: "Test",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := NewProduct(tt.id, tt.prodName, tt.price, tt.stock, tt.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && product == nil {
				t.Error("NewProduct() returned nil product without error")
			}
		})
	}
}

func TestProduct_ReduceStock(t *testing.T) {
	product, _ := NewProduct("PROD-001", "Test", 100, 10, "Test")

	tests := []struct {
		name     string
		quantity int
		wantErr  bool
		expected int
	}{
		{
			name:     "reduce within stock",
			quantity: 5,
			wantErr:  false,
			expected: 5,
		},
		{
			name:     "reduce exact stock",
			quantity: 5,
			wantErr:  false,
			expected: 0,
		},
		{
			name:     "reduce more than stock",
			quantity: 1,
			wantErr:  true,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := product.ReduceStock(tt.quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReduceStock() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && product.Stock != tt.expected {
				t.Errorf("ReduceStock() stock = %v, want %v", product.Stock, tt.expected)
			}
		})
	}
}