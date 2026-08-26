package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"order-service/internal/model"
	"order-service/internal/repository"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetOverviewToday(ctx context.Context, restID uuid.UUID) (*model.DashboardOverview, error) {
	return s.repo.GetOverviewToday(ctx, restID)
}

func (s *DashboardService) GetActiveOrders(ctx context.Context, restID uuid.UUID) ([]model.Order, error) {
	return s.repo.GetActiveOrders(ctx, restID)
}

func (s *DashboardService) GetRevenueStats(ctx context.Context, restID uuid.UUID, from, to time.Time, groupBy string) ([]model.RevenueStat, error) {
	if groupBy == "" {
		groupBy = "day"
	}
	return s.repo.GetRevenueStats(ctx, restID, from, to, groupBy)
}

func (s *DashboardService) GetTopItems(ctx context.Context, restID uuid.UUID) ([]model.TopItemStat, error) {
	return s.repo.GetTopItems(ctx, restID, 10)
}

func (s *DashboardService) GetOrderCounts(ctx context.Context, restID uuid.UUID) ([]model.OrderCountStat, error) {
	return s.repo.GetOrderCounts(ctx, restID)
}
