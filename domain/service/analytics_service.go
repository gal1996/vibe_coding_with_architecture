package service

import (
	"context"
	"sort"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
)

// AnalyticsService handles analytics and reporting logic
type AnalyticsService struct {
	orderRepo     repository.OrderRepository
	productRepo   repository.ProductRepository
	stockRepo     repository.StockRepository
	warehouseRepo repository.WarehouseRepository
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	stockRepo repository.StockRepository,
	warehouseRepo repository.WarehouseRepository,
) *AnalyticsService {
	return &AnalyticsService{
		orderRepo:     orderRepo,
		productRepo:   productRepo,
		stockRepo:     stockRepo,
		warehouseRepo: warehouseRepo,
	}
}

// SalesReport represents the sales analytics report
type SalesReport struct {
	SalesSummary     SalesSummary           `json:"sales_summary"`
	TopProducts      []ProductRanking       `json:"top_products"`
	WarehouseStock   []WarehouseStockStatus `json:"warehouse_stock"`
	CouponAnalytics  CouponAnalytics        `json:"coupon_analytics"`
}

// SalesSummary represents sales summary information
type SalesSummary struct {
	TotalRevenue int `json:"total_revenue"`
	TotalOrders  int `json:"total_orders"`
}

// ProductRanking represents product sales ranking
type ProductRanking struct {
	ProductName string `json:"product_name"`
	TotalSold   int    `json:"total_sold"`
}

// WarehouseStockStatus represents stock status per warehouse
type WarehouseStockStatus struct {
	WarehouseID   string `json:"warehouse_id"`
	WarehouseName string `json:"warehouse_name"`
	TotalStock    int    `json:"total_stock"`
}

// CouponAnalytics represents coupon usage analytics
type CouponAnalytics struct {
	CouponUsageRate float64 `json:"coupon_usage_rate"` // percentage
}

// GenerateSalesReport generates a comprehensive sales report
func (s *AnalyticsService) GenerateSalesReport(ctx context.Context) (*SalesReport, error) {
	// Get all orders
	allOrders, err := s.orderRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate sales summary
	salesSummary := s.calculateSalesSummary(allOrders)

	// Calculate top products
	topProducts := s.calculateTopProducts(allOrders)

	// Calculate warehouse stock summary
	warehouseStock, err := s.calculateWarehouseStock(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate coupon analytics
	couponAnalytics := s.calculateCouponAnalytics(allOrders)

	return &SalesReport{
		SalesSummary:     salesSummary,
		TopProducts:      topProducts,
		WarehouseStock:   warehouseStock,
		CouponAnalytics:  couponAnalytics,
	}, nil
}

// calculateSalesSummary calculates total revenue and orders from completed orders
func (s *AnalyticsService) calculateSalesSummary(orders []*entity.Order) SalesSummary {
	totalRevenue := 0
	totalOrders := 0

	for _, order := range orders {
		if order.Status == entity.OrderStatusCompleted {
			totalRevenue += order.TotalPrice
			totalOrders++
		}
	}

	return SalesSummary{
		TotalRevenue: totalRevenue,
		TotalOrders:  totalOrders,
	}
}

// calculateTopProducts calculates the top 3 products by quantity sold
func (s *AnalyticsService) calculateTopProducts(orders []*entity.Order) []ProductRanking {
	productSales := make(map[string]struct {
		name     string
		quantity int
	})

	// Aggregate product sales from completed orders
	for _, order := range orders {
		if order.Status == entity.OrderStatusCompleted {
			for _, item := range order.Items {
				if existing, ok := productSales[item.ProductID]; ok {
					productSales[item.ProductID] = struct {
						name     string
						quantity int
					}{
						name:     existing.name,
						quantity: existing.quantity + item.Quantity,
					}
				} else {
					productSales[item.ProductID] = struct {
						name     string
						quantity int
					}{
						name:     item.ProductName,
						quantity: item.Quantity,
					}
				}
			}
		}
	}

	// Convert to slice for sorting
	var rankings []ProductRanking
	for _, product := range productSales {
		rankings = append(rankings, ProductRanking{
			ProductName: product.name,
			TotalSold:   product.quantity,
		})
	}

	// Sort by quantity sold (descending)
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].TotalSold > rankings[j].TotalSold
	})

	// Return top 3
	if len(rankings) > 3 {
		return rankings[:3]
	}
	return rankings
}

// calculateWarehouseStock calculates total stock per warehouse
func (s *AnalyticsService) calculateWarehouseStock(ctx context.Context) ([]WarehouseStockStatus, error) {
	// Get all warehouses
	warehouses, err := s.warehouseRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var warehouseStockList []WarehouseStockStatus

	for _, warehouse := range warehouses {
		// Get all stocks for this warehouse
		stocks, err := s.stockRepo.FindByWarehouseID(ctx, warehouse.ID)
		if err != nil {
			return nil, err
		}

		totalStock := 0
		for _, stock := range stocks {
			totalStock += stock.Quantity
		}

		warehouseStockList = append(warehouseStockList, WarehouseStockStatus{
			WarehouseID:   warehouse.ID,
			WarehouseName: warehouse.Name,
			TotalStock:    totalStock,
		})
	}

	return warehouseStockList, nil
}

// calculateCouponAnalytics calculates coupon usage statistics
func (s *AnalyticsService) calculateCouponAnalytics(orders []*entity.Order) CouponAnalytics {
	totalOrders := 0
	ordersWithCoupon := 0

	for _, order := range orders {
		// Only count completed orders (exclude payment failed)
		if order.Status == entity.OrderStatusCompleted {
			totalOrders++
			if order.AppliedCoupon != "" {
				ordersWithCoupon++
			}
		}
	}

	usageRate := 0.0
	if totalOrders > 0 {
		usageRate = float64(ordersWithCoupon) / float64(totalOrders) * 100
	}

	return CouponAnalytics{
		CouponUsageRate: usageRate,
	}
}