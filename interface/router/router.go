package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gal1996/vibe_coding_with_architecture/di"
)

// NewRouter creates and configures the router
func NewRouter(container *di.Container) *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		public := v1.Group("")
		{
			// User routes
			public.POST("/register", container.UserHandler.Register)
			public.POST("/login", container.UserHandler.Login)

			// Product routes (read-only for public)
			public.GET("/products", container.ProductHandler.ListProducts)
			public.GET("/products/:id", container.ProductHandler.GetProduct)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(container.AuthMiddleware.Authenticate())
		{
			// User profile
			protected.GET("/users/:id", container.UserHandler.GetProfile)

			// Order routes
			protected.POST("/orders", container.OrderHandler.CreateOrder)
			protected.GET("/orders", container.OrderHandler.ListUserOrders)
			protected.GET("/orders/:id", container.OrderHandler.GetOrder)

			// Admin-only routes
			admin := protected.Group("")
			admin.Use(container.AuthMiddleware.RequireAdmin())
			{
				// Product management
				admin.POST("/products", container.ProductHandler.CreateProduct)
			}
		}
	}

	return router
}