package di

import (
	"context"
	"fmt"

	"github.com/gal1996/vibe_coding_with_architecture/domain/entity"
	"github.com/gal1996/vibe_coding_with_architecture/domain/repository"
	"github.com/gal1996/vibe_coding_with_architecture/domain/service"
	"github.com/gal1996/vibe_coding_with_architecture/infrastructure/auth"
	"github.com/gal1996/vibe_coding_with_architecture/infrastructure/payment"
	"github.com/gal1996/vibe_coding_with_architecture/infrastructure/persistence"
	"github.com/gal1996/vibe_coding_with_architecture/interface/handler"
	"github.com/gal1996/vibe_coding_with_architecture/interface/middleware"
	"github.com/gal1996/vibe_coding_with_architecture/usecase/interactor"
	"github.com/gal1996/vibe_coding_with_architecture/usecase/port"
)

// Container holds all dependencies
type Container struct {
	// Repositories
	ProductRepository   repository.ProductRepository
	UserRepository      repository.UserRepository
	OrderRepository     repository.OrderRepository
	StockRepository     repository.StockRepository
	WarehouseRepository repository.WarehouseRepository

	// Services
	AuthService    port.AuthService
	PaymentService port.PaymentService
	OrderService   *service.OrderService
	StockService   *service.StockService

	// Use Cases
	ProductUseCase *interactor.ProductUseCase
	UserUseCase    *interactor.UserUseCase
	OrderUseCase   *interactor.OrderUseCase

	// Handlers
	ProductHandler *handler.ProductHandler
	UserHandler    *handler.UserHandler
	OrderHandler   *handler.OrderHandler

	// Middleware
	AuthMiddleware *middleware.AuthMiddleware
}

// NewContainer creates a new dependency injection container
func NewContainer() *Container {
	// Initialize repositories
	productRepo := persistence.NewMemoryProductRepository()
	userRepo := persistence.NewMemoryUserRepository()
	orderRepo := persistence.NewMemoryOrderRepository()
	stockRepo := persistence.NewMemoryStockRepository()
	warehouseRepo := persistence.NewMemoryWarehouseRepository()

	// Initialize services
	authService := auth.NewJWTAuthService(userRepo)
	paymentService := payment.NewSimulatedPaymentService()
	stockService := service.NewStockService(stockRepo, warehouseRepo)
	orderService := service.NewOrderService(productRepo, orderRepo, stockService)

	// Initialize use cases
	productUseCase := interactor.NewProductUseCase(productRepo, userRepo, authService, stockService)
	userUseCase := interactor.NewUserUseCase(userRepo, authService)
	orderUseCase := interactor.NewOrderUseCase(orderRepo, productRepo, orderService, authService, paymentService)

	// Initialize handlers
	productHandler := handler.NewProductHandler(productUseCase)
	userHandler := handler.NewUserHandler(userUseCase)
	orderHandler := handler.NewOrderHandler(orderUseCase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService, userRepo)

	return &Container{
		// Repositories
		ProductRepository:   productRepo,
		UserRepository:      userRepo,
		OrderRepository:     orderRepo,
		StockRepository:     stockRepo,
		WarehouseRepository: warehouseRepo,

		// Services
		AuthService:    authService,
		PaymentService: paymentService,
		OrderService:   orderService,
		StockService:   stockService,

		// Use Cases
		ProductUseCase: productUseCase,
		UserUseCase:    userUseCase,
		OrderUseCase:   orderUseCase,

		// Handlers
		ProductHandler: productHandler,
		UserHandler:    userHandler,
		OrderHandler:   orderHandler,

		// Middleware
		AuthMiddleware: authMiddleware,
	}
}

// SeedTestData seeds the container with test data
func (c *Container) SeedTestData() error {
	// Try to create admin user or get existing one
	adminUser, err := c.UserUseCase.Register(nil, interactor.RegisterInput{
		Username: "admin",
		Password: "admin123",
		IsAdmin:  true,
	})
	if err != nil {
		// If user already exists, try to login to get the user
		loginOutput, loginErr := c.UserUseCase.Login(nil, interactor.LoginInput{
			Username: "admin",
			Password: "admin123",
		})
		if loginErr != nil {
			return fmt.Errorf("failed to create or login admin user: %v", err)
		}
		adminUser = loginOutput.User
	}

	// Try to create regular user (ignore if exists)
	_, err = c.UserUseCase.Register(nil, interactor.RegisterInput{
		Username: "user",
		Password: "user123",
		IsAdmin:  false,
	})
	// Ignore error if user already exists

	// Create warehouses
	warehouse1, err := entity.NewWarehouse("WH-001", "東京倉庫", "東京都港区")
	if err != nil {
		return err
	}
	warehouse2, err := entity.NewWarehouse("WH-002", "大阪倉庫", "大阪府大阪市")
	if err != nil {
		return err
	}
	warehouse3, err := entity.NewWarehouse("WH-003", "福岡倉庫", "福岡県福岡市")
	if err != nil {
		return err
	}

	// Save warehouses
	err = c.WarehouseRepository.Create(nil, warehouse1)
	if err != nil {
		return err
	}
	err = c.WarehouseRepository.Create(nil, warehouse2)
	if err != nil {
		return err
	}
	err = c.WarehouseRepository.Create(nil, warehouse3)
	if err != nil {
		return err
	}

	// Create some products (using admin context)
	ctx := auth.SetUserInContext(context.Background(), adminUser)

	products := []struct {
		input interactor.CreateProductInput
		stocks map[string]int // warehouseID -> quantity
	}{
		{
			input: interactor.CreateProductInput{Name: "Laptop", Price: 1200, Category: "Electronics"},
			stocks: map[string]int{"WH-001": 5, "WH-002": 3, "WH-003": 2},
		},
		{
			input: interactor.CreateProductInput{Name: "Mouse", Price: 25, Category: "Electronics"},
			stocks: map[string]int{"WH-001": 20, "WH-002": 15, "WH-003": 15},
		},
		{
			input: interactor.CreateProductInput{Name: "Keyboard", Price: 75, Category: "Electronics"},
			stocks: map[string]int{"WH-001": 10, "WH-002": 10, "WH-003": 10},
		},
		{
			input: interactor.CreateProductInput{Name: "Desk", Price: 300, Category: "Furniture"},
			stocks: map[string]int{"WH-001": 2, "WH-002": 2, "WH-003": 1},
		},
		{
			input: interactor.CreateProductInput{Name: "Chair", Price: 150, Category: "Furniture"},
			stocks: map[string]int{"WH-001": 5, "WH-002": 5, "WH-003": 5},
		},
		{
			input: interactor.CreateProductInput{Name: "Coffee", Price: 10, Category: "Food"},
			stocks: map[string]int{"WH-001": 40, "WH-002": 30, "WH-003": 30},
		},
	}

	for _, p := range products {
		product, err := c.ProductUseCase.CreateProduct(ctx, p.input)
		if err != nil {
			return err
		}

		// Create stock entries for each warehouse
		for warehouseID, quantity := range p.stocks {
			stockID := fmt.Sprintf("STK-%s-%s", product.ID, warehouseID)
			stock, err := entity.NewStock(stockID, product.ID, warehouseID, quantity)
			if err != nil {
				return err
			}
			err = c.StockRepository.Create(nil, stock)
			if err != nil {
				return err
			}
		}
	}

	return nil
}