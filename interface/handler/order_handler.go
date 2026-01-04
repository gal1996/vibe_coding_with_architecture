package handler

import (
	"net/http"

	"github.com/gal1996/vibe_coding_with_architecture/usecase/interactor"
	"github.com/gin-gonic/gin"
)

// OrderHandler handles HTTP requests for orders
type OrderHandler struct {
	orderUseCase *interactor.OrderUseCase
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderUseCase *interactor.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
	}
}

// CreateOrderRequest represents the request body for creating an order
type CreateOrderRequest struct {
	Items []OrderItemRequest `json:"items" binding:"required,min=1"`
}

// OrderItemRequest represents an item in an order request
type OrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// CreateOrder handles POST /orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert request to use case input
	items := make([]interactor.OrderItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = interactor.OrderItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	input := interactor.CreateOrderInput{
		Items: items,
	}

	order, err := h.orderUseCase.CreateOrder(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetOrder handles GET /orders/:id
func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")

	order, err := h.orderUseCase.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// ListUserOrders handles GET /orders
func (h *OrderHandler) ListUserOrders(c *gin.Context) {
	orders, err := h.orderUseCase.ListUserOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"count":  len(orders),
	})
}
