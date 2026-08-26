package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"user-service/internal/model"
)

type TokenRepository struct {
	db *sqlx.DB
}

func NewTokenRepository(db *sqlx.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Save(ctx context.Context, token *model.RefreshToken) error {
	query := `
		INSERT INTO users.refresh_tokens (id, user_id, token_hash, expires_at)
		VALUES (:id, :user_id, :token_hash, :expires_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("TokenRepository.Save: %w", err)
	}
	return nil
}

// FindByHash tìm refresh token còn hiệu lực theo hash
func (r *TokenRepository) FindByHash(ctx context.Context, hash string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := r.db.GetContext(ctx, &token, `
		SELECT * FROM users.refresh_tokens
		WHERE token_hash = $1
		  AND revoked_at IS NULL
		  AND expires_at > NOW()
	`, hash)
	if err != nil {
		return nil, fmt.Errorf("TokenRepository.FindByHash: %w", err)
	}
	return &token, nil
}

// Revoke đánh dấu 1 token đã bị thu hồi
func (r *TokenRepository) Revoke(ctx context.Context, hash string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users.refresh_tokens
		SET revoked_at = NOW()
		WHERE token_hash = $1
	`, hash)
	return err
}

// RevokeAllByUser thu hồi tất cả token của 1 user (dùng khi đổi mật khẩu)
func (r *TokenRepository) RevokeAllByUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users.refresh_tokens
		SET revoked_at = NOW()
		WHERE user_id = $1 AND revoked_at IS NULL
	`, userID)
	return err
}
