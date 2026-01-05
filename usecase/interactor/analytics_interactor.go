package interactor

import (
	"context"
	"fmt"

	"github.com/gal1996/vibe_coding_with_architecture/domain/service"
	"github.com/gal1996/vibe_coding_with_architecture/usecase/port"
)

// AnalyticsUseCase handles analytics-related use cases
type AnalyticsUseCase struct {
	analyticsService *service.AnalyticsService
	authService      port.AuthService
}

// NewAnalyticsUseCase creates a new analytics use case
func NewAnalyticsUseCase(
	analyticsService *service.AnalyticsService,
	authService port.AuthService,
) *AnalyticsUseCase {
	return &AnalyticsUseCase{
		analyticsService: analyticsService,
		authService:      authService,
	}
}

// GetSalesReport generates a sales report (admin only)
func (uc *AnalyticsUseCase) GetSalesReport(ctx context.Context) (*service.SalesReport, error) {
	// Check if current user is admin
	currentUser, err := uc.authService.GetCurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}

	if !currentUser.IsAdmin {
		return nil, fmt.Errorf("permission denied: admin access required")
	}

	// Generate report using analytics service
	report, err := uc.analyticsService.GenerateSalesReport(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate sales report: %w", err)
	}

	return report, nil
}