package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ─── Types ────────────────────────────────────────────────────────────────────

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpire  time.Duration
	RefreshExpire time.Duration
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string    // raw token (gửi cho client)
	RefreshHash  string    // hash (lưu DB)
	ExpiresAt    time.Time // refresh token expires
}

// ─── Generate ─────────────────────────────────────────────────────────────────

// GenerateTokenPair tạo JWT access token và random refresh token
func GenerateTokenPair(userID uuid.UUID, role string, cfg JWTConfig) (*TokenPair, error) {
	// Access Token — JWT signed HS256
	accessClaims := JWTClaims{
		UserID: userID.String(),
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.AccessExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}
	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessJWT.SignedString([]byte(cfg.AccessSecret))
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	// Refresh Token — 32 random bytes, lưu SHA-256 hash trong DB
	rawBytes := make([]byte, 32)
	if _, err := rand.Read(rawBytes); err != nil {
		return nil, fmt.Errorf("generate refresh token bytes: %w", err)
	}
	rawToken := hex.EncodeToString(rawBytes)
	hash := sha256.Sum256([]byte(rawToken))
	refreshHash := hex.EncodeToString(hash[:])

	return &TokenPair{
		AccessToken:  accessTokenStr,
		RefreshToken: rawToken,
		RefreshHash:  refreshHash,
		ExpiresAt:    time.Now().Add(cfg.RefreshExpire),
	}, nil
}

// ─── Validate ─────────────────────────────────────────────────────────────────

// ValidateAccessToken parse và validate JWT access token
func ValidateAccessToken(tokenStr, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

// HashRefreshToken tính SHA-256 hash của raw refresh token
func HashRefreshToken(rawToken string) string {
	h := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(h[:])
}
