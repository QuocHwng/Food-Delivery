package model

import (
	"time"

	"github.com/google/uuid"
)

// ─── Database Models ─────────────────────────────────────────────────────────

type Restaurant struct {
	ID             uuid.UUID `db:"id" json:"id"`
	OwnerID        uuid.UUID `db:"owner_id" json:"owner_id"`
	Name           string    `db:"name" json:"name"`
	Description    *string   `db:"description" json:"description"`
	Address        *string   `db:"address" json:"address"`
	Phone          *string   `db:"phone" json:"phone"`
	AvgRating      float64   `db:"avg_rating" json:"avg_rating"`
	TotalRatings   int       `db:"total_ratings" json:"total_ratings"`
	MinOrderValue  float64   `db:"min_order_value" json:"min_order_value"`
	DeliveryFee    float64   `db:"delivery_fee" json:"delivery_fee"`
	IsOpen         bool      `db:"is_open" json:"is_open"`
	IsActive       bool      `db:"is_active" json:"is_active"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

type MenuCategory struct {
	ID           uuid.UUID `db:"id" json:"id"`
	RestaurantID uuid.UUID `db:"restaurant_id" json:"restaurant_id"`
	Name         string    `db:"name" json:"name"`
	DisplayOrder int       `db:"display_order" json:"display_order"`
	IsActive     bool      `db:"is_active" json:"is_active"`
}

type MenuItem struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	RestaurantID uuid.UUID  `db:"restaurant_id" json:"restaurant_id"`
	CategoryID   *uuid.UUID `db:"category_id" json:"category_id"`
	Name         string     `db:"name" json:"name"`
	Description  *string    `db:"description" json:"description"`
	Price        float64    `db:"price" json:"price"`
	ImageURL     *string    `db:"image_url" json:"image_url"`
	IsAvailable  bool       `db:"is_available" json:"is_available"`
}

// ─── Requests & Responses ────────────────────────────────────────────────────

type CreateRestaurantReq struct {
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description"`
	Address       string  `json:"address" binding:"required"`
	Phone         string  `json:"phone"`
	MinOrderValue float64 `json:"min_order_value"`
	DeliveryFee   float64 `json:"delivery_fee"`
}

type CreateMenuCategoryReq struct {
	Name         string `json:"name" binding:"required"`
	DisplayOrder int    `json:"display_order"`
}

type CreateMenuItemReq struct {
	CategoryID  *uuid.UUID `json:"category_id"`
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	Price       float64    `json:"price" binding:"required,min=0"`
	ImageURL    string     `json:"image_url"`
}

// MenuResponse groups categories with their items for a complete menu view
type MenuResponse struct {
	Categories []CategoryWithItems `json:"categories"`
	Uncategorized []MenuItem       `json:"uncategorized"`
}

type CategoryWithItems struct {
	MenuCategory
	Items []MenuItem `json:"items"`
}
