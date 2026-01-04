package repository

import (
	"context"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
)

// StockRepository defines the interface for stock persistence
type StockRepository interface {
	Create(ctx context.Context, stock *entity.Stock) error
	Update(ctx context.Context, stock *entity.Stock) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*entity.Stock, error)
	FindByProductID(ctx context.Context, productID string) ([]*entity.Stock, error)
	FindByWarehouseID(ctx context.Context, warehouseID string) ([]*entity.Stock, error)
	FindByProductAndWarehouse(ctx context.Context, productID, warehouseID string) (*entity.Stock, error)

	// Transaction support
	BeginTransaction(ctx context.Context) (StockTransaction, error)
}

// StockTransaction represents a stock transaction
type StockTransaction interface {
	GetStockRepository() StockRepository
	Commit() error
	Rollback() error
}