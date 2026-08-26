package model

import "time"

type DashboardOverview struct {
	OrdersToday    int     `json:"orders_today"`
	RevenueToday   float64 `json:"revenue_today"`
	PendingOrders  int     `json:"pending_orders"`
}

type RevenueStat struct {
	Period  time.Time `db:"period" json:"period"`
	Revenue float64   `db:"revenue" json:"revenue"`
}

type TopItemStat struct {
	MenuItemID string  `db:"menu_item_id" json:"menu_item_id"`
	Name       string  `db:"name" json:"name"`
	TotalSold  int     `db:"total_sold" json:"total_sold"`
}

type OrderCountStat struct {
	Status string `db:"status" json:"status"`
	Count  int    `db:"count" json:"count"`
}
