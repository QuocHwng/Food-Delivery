package model

import (
	"time"

	"github.com/google/uuid"
)

// ─── DB Model ─────────────────────────────────────────────────────────────────

type UserAddress struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	UserID    uuid.UUID `db:"user_id"    json:"user_id"`
	Label     *string   `db:"label"      json:"label,omitempty"`
	Street    string    `db:"street"     json:"street"`
	District  *string   `db:"district"   json:"district,omitempty"`
	City      *string   `db:"city"       json:"city,omitempty"`
	Lat       *float64  `db:"lat"        json:"lat,omitempty"`
	Lng       *float64  `db:"lng"        json:"lng,omitempty"`
	IsDefault bool      `db:"is_default" json:"is_default"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// ─── Requests ─────────────────────────────────────────────────────────────────

type CreateAddressRequest struct {
	Label    string   `json:"label"`
	Street   string   `json:"street"   binding:"required"`
	District string   `json:"district"`
	City     string   `json:"city"`
	Lat      *float64 `json:"lat"`
	Lng      *float64 `json:"lng"`
}

type UpdateAddressRequest struct {
	Label    string   `json:"label"`
	Street   string   `json:"street"`
	District string   `json:"district"`
	City     string   `json:"city"`
	Lat      *float64 `json:"lat"`
	Lng      *float64 `json:"lng"`
}

// ─── DB Model ─────────────────────────────────────────────────────────────────

type RefreshToken struct {
	ID        uuid.UUID  `db:"id"`
	UserID    uuid.UUID  `db:"user_id"`
	TokenHash string     `db:"token_hash"`
	ExpiresAt time.Time  `db:"expires_at"`
	RevokedAt *time.Time `db:"revoked_at"`
	CreatedAt time.Time  `db:"created_at"`
}
