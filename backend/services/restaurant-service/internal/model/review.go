package model

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	RestaurantID  uuid.UUID  `db:"restaurant_id" json:"restaurant_id"`
	OrderID       uuid.UUID  `db:"order_id" json:"order_id"`
	CustomerID    uuid.UUID  `db:"customer_id" json:"customer_id"`
	Score         int        `db:"score" json:"score"`
	Comment       *string    `db:"comment" json:"comment"`
	ImageURL      *string    `db:"image_url" json:"image_url"`
	OwnerReply    *string    `db:"owner_reply" json:"owner_reply"`
	OwnerReplyAt  *time.Time `db:"owner_reply_at" json:"owner_reply_at"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`

	Items []ReviewItem `db:"-" json:"items,omitempty"`
}

type ReviewItem struct {
	ID         uuid.UUID `db:"id" json:"id"`
	ReviewID   uuid.UUID `db:"review_id" json:"review_id"`
	MenuItemID uuid.UUID `db:"menu_item_id" json:"menu_item_id"`
	Score      int       `db:"score" json:"score"`
}

// Request payloads
type CreateReviewReq struct {
	RestaurantID uuid.UUID            `json:"restaurant_id" binding:"required"`
	Score        int                  `json:"score" binding:"required,min=1,max=5"`
	Comment      string               `json:"comment"`
	ImageURL     string               `json:"image_url"`
	Items        []CreateReviewItemReq `json:"items"`
}

type CreateReviewItemReq struct {
	MenuItemID uuid.UUID `json:"menu_item_id" binding:"required"`
	Score      int       `json:"score" binding:"required,min=1,max=5"`
}

type ReplyReviewReq struct {
	Reply string `json:"reply" binding:"required"`
}
