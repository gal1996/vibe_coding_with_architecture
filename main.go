package main

import (
	"log"
	"os"

	"github.com/gal1996/vibe_coding_with_architecture/di"
	"github.com/gal1996/vibe_coding_with_architecture/interface/router"
)

func main() {
	// Initialize dependency injection container
	container := di.NewContainer()

	// Seed test data
	if err := container.SeedTestData(); err != nil {
		log.Printf("Warning: Failed to seed test data: %v", err)
	} else {
		log.Println("Test data seeded successfully")
		log.Println("Available test accounts:")
		log.Println("  Admin: username=admin, password=admin123")
		log.Println("  User:  username=user, password=user123")
	}

	// Create router
	r := router.NewRouter(container)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Starting EC Site Backend API on port %s", port)
	log.Printf("API Documentation:")
	log.Printf("  Health Check: GET /health")
	log.Printf("  Register:     POST /api/v1/register")
	log.Printf("  Login:        POST /api/v1/login")
	log.Printf("  Products:     GET /api/v1/products")
	log.Printf("  Create Order: POST /api/v1/orders (requires auth)")
	log.Printf("  Create Product: POST /api/v1/products (requires admin)")

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}