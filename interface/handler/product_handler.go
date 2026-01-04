package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gal1996/vibe_coding_with_architecture/usecase/interactor"
)

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
	productUseCase *interactor.ProductUseCase
}

// NewProductHandler creates a new product handler
func NewProductHandler(productUseCase *interactor.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
	}
}

// CreateProductRequest represents the request body for creating a product
type CreateProductRequest struct {
	Name     string `json:"name" binding:"required"`
	Price    int    `json:"price" binding:"required,min=0"`
	Stock    int    `json:"stock" binding:"required,min=0"`
	Category string `json:"category" binding:"required"`
}

// CreateProduct handles POST /products
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := interactor.CreateProductInput{
		Name:     req.Name,
		Price:    req.Price,
		Stock:    req.Stock,
		Category: req.Category,
	}

	product, err := h.productUseCase.CreateProduct(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct handles GET /products/:id
func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID := c.Param("id")

	product, err := h.productUseCase.GetProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// ListProducts handles GET /products
func (h *ProductHandler) ListProducts(c *gin.Context) {
	category := c.Query("category")

	products, err := h.productUseCase.ListProducts(c.Request.Context(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"count":    len(products),
	})
}