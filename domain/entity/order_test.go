package entity

import (
	"testing"
)

func TestOrder_CalculateTotalWithTaxAndShipping(t *testing.T) {
	tests := []struct {
		name            string
		items           []struct {
			productID   string
			productName string
			quantity    int
			price       int
		}
		expectedSubtotal    int
		expectedTax         int
		expectedShippingFee int
		expectedTotal       int
	}{
		{
			name: "Order under 5000 yen - should add 500 yen shipping",
			items: []struct {
				productID   string
				productName string
				quantity    int
				price       int
			}{
				{"P001", "Product 1", 2, 1000}, // 2000 yen
				{"P002", "Product 2", 1, 2000}, // 2000 yen
			},
			expectedSubtotal:    4000,
			expectedTax:         400,  // 10% of 4000
			expectedShippingFee: 500,  // Under 5000 yen
			expectedTotal:       4900, // 4000 + 400 + 500
		},
		{
			name: "Order exactly 5000 yen - free shipping",
			items: []struct {
				productID   string
				productName string
				quantity    int
				price       int
			}{
				{"P001", "Product 1", 5, 1000}, // 5000 yen
			},
			expectedSubtotal:    5000,
			expectedTax:         500,  // 10% of 5000
			expectedShippingFee: 0,    // Free shipping
			expectedTotal:       5500, // 5000 + 500 + 0
		},
		{
			name: "Order over 5000 yen - free shipping",
			items: []struct {
				productID   string
				productName string
				quantity    int
				price       int
			}{
				{"P001", "Product 1", 3, 2000}, // 6000 yen
				{"P002", "Product 2", 2, 1500}, // 3000 yen
			},
			expectedSubtotal:    9000,
			expectedTax:         900,   // 10% of 9000
			expectedShippingFee: 0,     // Free shipping
			expectedTotal:       9900,  // 9000 + 900 + 0
		},
		{
			name: "Small order - should add shipping",
			items: []struct {
				productID   string
				productName string
				quantity    int
				price       int
			}{
				{"P001", "Product 1", 1, 500}, // 500 yen
			},
			expectedSubtotal:    500,
			expectedTax:         50,   // 10% of 500
			expectedShippingFee: 500,  // Under 5000 yen
			expectedTotal:       1050, // 500 + 50 + 500
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create order
			order, err := NewOrder("ORD-001", "USER-001")
			if err != nil {
				t.Fatalf("Failed to create order: %v", err)
			}

			// Add items
			for _, item := range tt.items {
				err = order.AddItem(item.productID, item.productName, item.quantity, item.price)
				if err != nil {
					t.Fatalf("Failed to add item: %v", err)
				}
			}

			// Check subtotal
			if subtotal := order.GetSubtotal(); subtotal != tt.expectedSubtotal {
				t.Errorf("Expected subtotal %d, got %d", tt.expectedSubtotal, subtotal)
			}

			// Check tax amount
			if taxAmount := order.GetTaxAmount(); taxAmount != tt.expectedTax {
				t.Errorf("Expected tax amount %d, got %d", tt.expectedTax, taxAmount)
			}

			// Check shipping fee
			if order.ShippingFee != tt.expectedShippingFee {
				t.Errorf("Expected shipping fee %d, got %d", tt.expectedShippingFee, order.ShippingFee)
			}

			// Check total price
			if order.TotalPrice != tt.expectedTotal {
				t.Errorf("Expected total price %d, got %d", tt.expectedTotal, order.TotalPrice)
			}
		})
	}
}

func TestOrder_ShippingFeeCalculation(t *testing.T) {
	order, _ := NewOrder("ORD-001", "USER-001")

	// Test with subtotal under 5000
	order.AddItem("P001", "Product 1", 2, 2000) // 4000 yen
	if fee := order.CalculateShippingFee(); fee != 500 {
		t.Errorf("Expected shipping fee 500 for order under 5000 yen, got %d", fee)
	}

	// Test with subtotal exactly 5000
	order.Items = []OrderItem{} // Reset items
	order.AddItem("P001", "Product 1", 5, 1000) // 5000 yen
	if fee := order.CalculateShippingFee(); fee != 0 {
		t.Errorf("Expected shipping fee 0 for order of 5000 yen or more, got %d", fee)
	}

	// Test with subtotal over 5000
	order.Items = []OrderItem{} // Reset items
	order.AddItem("P001", "Product 1", 3, 2000) // 6000 yen
	if fee := order.CalculateShippingFee(); fee != 0 {
		t.Errorf("Expected shipping fee 0 for order over 5000 yen, got %d", fee)
	}
}