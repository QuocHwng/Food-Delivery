package model

import (
	"time"

	"github.com/google/uuid"
)

// ─── DB Model ─────────────────────────────────────────────────────────────────

type User struct {
	ID           uuid.UUID `db:"id"            json:"id"`
	Name         string    `db:"name"          json:"name"`
	Email        string    `db:"email"         json:"email"`
	Phone        *string   `db:"phone"         json:"phone,omitempty"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         string    `db:"role"          json:"role"`
	AvatarURL    *string   `db:"avatar_url"    json:"avatar_url,omitempty"`
	IsVerified   bool      `db:"is_verified"   json:"is_verified"`
	IsActive     bool      `db:"is_active"     json:"is_active"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"    json:"updated_at"`
}

// ─── Requests ─────────────────────────────────────────────────────────────────

type RegisterRequest struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Phone    string `json:"phone"    binding:"omitempty,min=9,max=15"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"     binding:"required,oneof=customer restaurant_owner"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdateProfileRequest struct {
	Name      string `json:"name"       binding:"omitempty,min=2,max=100"`
	Phone     string `json:"phone"      binding:"omitempty"`
	AvatarURL string `json:"avatar_url" binding:"omitempty"`
}

// ─── Responses ────────────────────────────────────────────────────────────────

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}
