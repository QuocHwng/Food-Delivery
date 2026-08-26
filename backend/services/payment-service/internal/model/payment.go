package model

import (
	"time"
	"github.com/google/uuid"
)

type Payment struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	OrderID       uuid.UUID  `db:"order_id" json:"order_id"`
	Amount        float64    `db:"amount" json:"amount"`
	Method        string     `db:"method" json:"method"`
	Status        string     `db:"status" json:"status"`
	TransactionID *string    `db:"vnpay_txn_ref" json:"vnpay_txn_ref"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
}

type PaymentLog struct {
	ID        uuid.UUID `db:"id" json:"id"`
	PaymentID uuid.UUID `db:"payment_id" json:"payment_id"`
	Event     string    `db:"event" json:"event"`
	Direction string    `db:"direction" json:"direction"`
	Payload   string    `db:"payload" json:"payload"` // using JSON string for simplicity in raw SQL
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Requests
type CreatePaymentReq struct {
	OrderID    uuid.UUID `json:"order_id" binding:"required"`
	Amount     float64   `json:"amount" binding:"required,min=1000"`
	Method     string    `json:"method" binding:"required,oneof=vnpay cod"`
	IPAddress  string    `json:"-"`
}
