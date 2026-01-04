package repository

import (
	"context"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
)

// WarehouseRepository defines the interface for warehouse persistence
type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *entity.Warehouse) error
	Update(ctx context.Context, warehouse *entity.Warehouse) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*entity.Warehouse, error)
	FindAll(ctx context.Context) ([]*entity.Warehouse, error)
}