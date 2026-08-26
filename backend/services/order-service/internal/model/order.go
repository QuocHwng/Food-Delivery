package model

import (
	"time"
	"github.com/google/uuid"
)

// ─── Entities ─────────────────────────────────────────────────────────────────

type Order struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	CustomerID     uuid.UUID  `db:"customer_id" json:"customer_id"`
	RestaurantID   uuid.UUID  `db:"restaurant_id" json:"restaurant_id"`
	AddressID      *uuid.UUID `db:"address_id" json:"address_id"`
	Status         string     `db:"status" json:"status"`
	Subtotal       float64    `db:"subtotal" json:"subtotal"`
	DiscountAmount float64    `db:"discount_amount" json:"discount_amount"`
	DeliveryFee    float64    `db:"delivery_fee" json:"delivery_fee"`
	Total          float64    `db:"total" json:"total"`
	Note           *string    `db:"note" json:"note"`
	CouponID       *uuid.UUID `db:"coupon_id" json:"coupon_id"`
	CancelledBy    *uuid.UUID `db:"cancelled_by" json:"cancelled_by"`
	CancelReason   *string    `db:"cancel_reason" json:"cancel_reason"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`

	Items []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID              uuid.UUID `db:"id" json:"id"`
	OrderID         uuid.UUID `db:"order_id" json:"order_id"`
	MenuItemID      uuid.UUID `db:"menu_item_id" json:"menu_item_id"`
	NameSnapshot    string    `db:"name_snapshot" json:"name_snapshot"`
	PriceSnapshot   float64   `db:"price_snapshot" json:"price_snapshot"`
	Quantity        int       `db:"quantity" json:"quantity"`
	OptionsSnapshot *string   `db:"options_snapshot" json:"options_snapshot"` // stored as JSONB string for now
}

type OrderStatusHistory struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	OrderID    uuid.UUID  `db:"order_id" json:"order_id"`
	FromStatus *string    `db:"from_status" json:"from_status"`
	ToStatus   string     `db:"to_status" json:"to_status"`
	ChangedBy  uuid.UUID  `db:"changed_by" json:"changed_by"`
	Note       *string    `db:"note" json:"note"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
}

// ─── Requests ─────────────────────────────────────────────────────────────────

type CreateOrderReq struct {
	RestaurantID uuid.UUID          `json:"restaurant_id" binding:"required"`
	AddressID    uuid.UUID          `json:"address_id" binding:"required"`
	Note         string             `json:"note"`
	Items        []CreateOrderItem  `json:"items" binding:"required,min=1"`
}

type CreateOrderItem struct {
	MenuItemID uuid.UUID `json:"menu_item_id" binding:"required"`
	Quantity   int       `json:"quantity" binding:"required,min=1"`
}

type UpdateStatusReq struct {
	Status string `json:"status" binding:"required,oneof=confirmed preparing ready completed"`
	Note   string `json:"note"`
}

type CancelOrderReq struct {
	Reason string `json:"reason" binding:"required"`
}
