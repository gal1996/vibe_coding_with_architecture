package handler

import (
	"net/http"

	"github.com/gal1996/vibe_coding_with_architecture/usecase/interactor"
	"github.com/gin-gonic/gin"
)

// AdminHandler handles admin-related HTTP requests
type AdminHandler struct {
	analyticsUseCase *interactor.AnalyticsUseCase
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(analyticsUseCase *interactor.AnalyticsUseCase) *AdminHandler {
	return &AdminHandler{
		analyticsUseCase: analyticsUseCase,
	}
}

// GetSalesReport handles GET /admin/reports/sales
func (h *AdminHandler) GetSalesReport(c *gin.Context) {
	report, err := h.analyticsUseCase.GetSalesReport(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}