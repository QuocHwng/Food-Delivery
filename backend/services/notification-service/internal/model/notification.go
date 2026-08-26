package model

import (
	"time"
	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	UserID    uuid.UUID  `db:"user_id" json:"user_id"`
	Title     string     `db:"title" json:"title"`
	Message   string     `db:"message" json:"message"`
	Type      string     `db:"type" json:"type"` // "order_status", "payment", "system"
	IsRead    bool       `db:"is_read" json:"is_read"`
	Data      *string    `db:"data" json:"data"` // JSON string for extra payload (e.g. order_id)
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

// RabbitMQ Event Payload
type NotificationEvent struct {
	UserID  uuid.UUID              `json:"user_id"`
	Title   string                 `json:"title"`
	Message string                 `json:"message"`
	Type    string                 `json:"type"`
	Data    map[string]interface{} `json:"data"`
}
