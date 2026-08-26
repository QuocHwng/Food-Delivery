package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"order-service/internal/model"
)

type DashboardRepository struct {
	db *sqlx.DB
}

func NewDashboardRepository(db *sqlx.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetOverviewToday(ctx context.Context, restID uuid.UUID) (*model.DashboardOverview, error) {
	var overview model.DashboardOverview
	query := `
		SELECT 
			COUNT(id) FILTER (WHERE status = 'completed') as orders_today,
			COALESCE(SUM(total) FILTER (WHERE status = 'completed'), 0) as revenue_today,
			COUNT(id) FILTER (WHERE status = 'pending') as pending_orders
		FROM orders.orders 
		WHERE restaurant_id = $1 AND DATE(created_at AT TIME ZONE 'UTC') = DATE(NOW() AT TIME ZONE 'UTC')
	`
	err := r.db.GetContext(ctx, &overview, query, restID)
	return &overview, err
}

func (r *DashboardRepository) GetActiveOrders(ctx context.Context, restID uuid.UUID) ([]model.Order, error) {
	var list []model.Order
	query := `SELECT * FROM orders.orders WHERE restaurant_id = $1 AND status IN ('pending', 'confirmed', 'preparing') ORDER BY created_at ASC`
	err := r.db.SelectContext(ctx, &list, query, restID)
	return list, err
}

func (r *DashboardRepository) GetRevenueStats(ctx context.Context, restID uuid.UUID, from, to time.Time, groupBy string) ([]model.RevenueStat, error) {
	var list []model.RevenueStat
	// Prevent SQL injection by strictly matching the group by parameter
	trunc := "day"
	if groupBy == "week" || groupBy == "month" {
		trunc = groupBy
	}
	
	query := `
		SELECT DATE_TRUNC($1, created_at) as period, SUM(total) as revenue
		FROM orders.orders
		WHERE restaurant_id = $2 AND status = 'completed' AND created_at >= $3 AND created_at <= $4
		GROUP BY period ORDER BY period
	`
	err := r.db.SelectContext(ctx, &list, query, trunc, restID, from, to)
	return list, err
}

func (r *DashboardRepository) GetTopItems(ctx context.Context, restID uuid.UUID, limit int) ([]model.TopItemStat, error) {
	var list []model.TopItemStat
	query := `
		SELECT oi.menu_item_id::text, oi.name_snapshot as name, SUM(oi.quantity) as total_sold
		FROM orders.order_items oi
		JOIN orders.orders o ON oi.order_id = o.id
		WHERE o.restaurant_id = $1 AND o.status = 'completed'
		GROUP BY oi.menu_item_id, oi.name_snapshot
		ORDER BY total_sold DESC LIMIT $2
	`
	err := r.db.SelectContext(ctx, &list, query, restID, limit)
	return list, err
}

func (r *DashboardRepository) GetOrderCounts(ctx context.Context, restID uuid.UUID) ([]model.OrderCountStat, error) {
	var list []model.OrderCountStat
	query := `
		SELECT status, COUNT(id) as count
		FROM orders.orders
		WHERE restaurant_id = $1
		GROUP BY status
	`
	err := r.db.SelectContext(ctx, &list, query, restID)
	return list, err
}
