package di

import (
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
	ProductRepository repository.ProductRepository
	UserRepository    repository.UserRepository
	OrderRepository   repository.OrderRepository

	// Services
	AuthService    port.AuthService
	PaymentService port.PaymentService
	OrderService   *service.OrderService

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

	// Initialize services
	authService := auth.NewJWTAuthService(userRepo)
	paymentService := payment.NewSimulatedPaymentService()
	orderService := service.NewOrderService(productRepo, orderRepo)

	// Initialize use cases
	productUseCase := interactor.NewProductUseCase(productRepo, userRepo, authService)
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
		ProductRepository: productRepo,
		UserRepository:    userRepo,
		OrderRepository:   orderRepo,

		// Services
		AuthService:    authService,
		PaymentService: paymentService,
		OrderService:   orderService,

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
	// Create admin user
	adminUser, err := c.UserUseCase.Register(nil, interactor.RegisterInput{
		Username: "admin",
		Password: "admin123",
		IsAdmin:  true,
	})
	if err != nil {
		return err
	}

	// Create regular user
	_, err = c.UserUseCase.Register(nil, interactor.RegisterInput{
		Username: "user",
		Password: "user123",
		IsAdmin:  false,
	})
	if err != nil {
		return err
	}

	// Create some products (using admin context)
	ctx := auth.SetUserInContext(nil, adminUser)

	products := []interactor.CreateProductInput{
		{Name: "Laptop", Price: 1200, Stock: 10, Category: "Electronics"},
		{Name: "Mouse", Price: 25, Stock: 50, Category: "Electronics"},
		{Name: "Keyboard", Price: 75, Stock: 30, Category: "Electronics"},
		{Name: "Desk", Price: 300, Stock: 5, Category: "Furniture"},
		{Name: "Chair", Price: 150, Stock: 15, Category: "Furniture"},
		{Name: "Coffee", Price: 10, Stock: 100, Category: "Food"},
	}

	for _, p := range products {
		_, err = c.ProductUseCase.CreateProduct(ctx, p)
		if err != nil {
			return err
		}
	}

	return nil
}