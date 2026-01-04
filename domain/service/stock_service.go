package service

import (
	"context"
	"fmt"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
)

// StockService handles stock management across warehouses
type StockService struct {
	stockRepo     repository.StockRepository
	warehouseRepo repository.WarehouseRepository
}

// NewStockService creates a new stock service
func NewStockService(stockRepo repository.StockRepository, warehouseRepo repository.WarehouseRepository) *StockService {
	return &StockService{
		stockRepo:     stockRepo,
		warehouseRepo: warehouseRepo,
	}
}

// StockAllocation represents how stock is allocated from different warehouses
type StockAllocation struct {
	WarehouseID   string
	WarehouseName string
	Quantity      int
}

// CheckAvailability checks if a product has sufficient stock across all warehouses
func (s *StockService) CheckAvailability(ctx context.Context, productID string, requiredQuantity int) (bool, int, error) {
	stocks, err := s.stockRepo.FindByProductID(ctx, productID)
	if err != nil {
		return false, 0, err
	}

	totalAvailable := 0
	for _, stock := range stocks {
		totalAvailable += stock.Quantity
	}

	return totalAvailable >= requiredQuantity, totalAvailable, nil
}

// AllocateStock allocates stock from multiple warehouses for an order
func (s *StockService) AllocateStock(ctx context.Context, productID string, requiredQuantity int) ([]StockAllocation, error) {
	// Start a transaction
	tx, err := s.stockRepo.BeginTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stockRepo := tx.GetStockRepository()

	// Get all stocks for the product
	stocks, err := stockRepo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to find stocks: %w", err)
	}

	// Check total availability
	totalAvailable := 0
	for _, stock := range stocks {
		totalAvailable += stock.Quantity
	}

	if totalAvailable < requiredQuantity {
		return nil, fmt.Errorf("insufficient stock: required=%d, available=%d", requiredQuantity, totalAvailable)
	}

	// Allocate stock from warehouses
	allocations := []StockAllocation{}
	remaining := requiredQuantity

	for _, stock := range stocks {
		if remaining == 0 {
			break
		}

		// Get warehouse details
		warehouse, err := s.warehouseRepo.FindByID(ctx, stock.WarehouseID)
		if err != nil {
			return nil, fmt.Errorf("failed to find warehouse %s: %w", stock.WarehouseID, err)
		}

		// Determine how much to allocate from this warehouse
		toAllocate := remaining
		if stock.Quantity < remaining {
			toAllocate = stock.Quantity
		}

		if toAllocate > 0 {
			// Reduce stock in this warehouse
			err = stock.Reduce(toAllocate)
			if err != nil {
				return nil, fmt.Errorf("failed to reduce stock: %w", err)
			}

			// Update the stock record
			err = stockRepo.Update(ctx, stock)
			if err != nil {
				return nil, fmt.Errorf("failed to update stock: %w", err)
			}

			allocations = append(allocations, StockAllocation{
				WarehouseID:   stock.WarehouseID,
				WarehouseName: warehouse.Name,
				Quantity:      toAllocate,
			})

			remaining -= toAllocate
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return allocations, nil
}

// RestoreStock restores stock to warehouses (for rollback scenarios)
func (s *StockService) RestoreStock(ctx context.Context, productID string, allocations []StockAllocation) error {
	for _, allocation := range allocations {
		stock, err := s.stockRepo.FindByProductAndWarehouse(ctx, productID, allocation.WarehouseID)
		if err != nil {
			// If stock record doesn't exist, create a new one
			stockID := fmt.Sprintf("STK-%s-%s", productID, allocation.WarehouseID)
			stock, err = entity.NewStock(stockID, productID, allocation.WarehouseID, 0)
			if err != nil {
				return fmt.Errorf("failed to create stock record: %w", err)
			}
			err = s.stockRepo.Create(ctx, stock)
			if err != nil {
				return fmt.Errorf("failed to create stock: %w", err)
			}
		}

		// Add back the allocated quantity
		err = stock.Add(allocation.Quantity)
		if err != nil {
			return fmt.Errorf("failed to restore stock: %w", err)
		}

		err = s.stockRepo.Update(ctx, stock)
		if err != nil {
			return fmt.Errorf("failed to update stock: %w", err)
		}
	}

	return nil
}

// GetProductStockInfo gets stock information for a product across all warehouses
func (s *StockService) GetProductStockInfo(ctx context.Context, productID string) ([]entity.StockInfo, int, error) {
	stocks, err := s.stockRepo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, 0, err
	}

	stockInfos := []entity.StockInfo{}
	totalStock := 0

	for _, stock := range stocks {
		warehouse, err := s.warehouseRepo.FindByID(ctx, stock.WarehouseID)
		if err != nil {
			continue // Skip if warehouse not found
		}

		stockInfos = append(stockInfos, entity.StockInfo{
			WarehouseID:   stock.WarehouseID,
			WarehouseName: warehouse.Name,
			Quantity:      stock.Quantity,
		})
		totalStock += stock.Quantity
	}

	return stockInfos, totalStock, nil
}