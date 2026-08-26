package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"user-service/internal/model"
)

type AddressRepository struct {
	db *sqlx.DB
}

func NewAddressRepository(db *sqlx.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

func (r *AddressRepository) Create(ctx context.Context, addr *model.UserAddress) error {
	query := `
		INSERT INTO users.user_addresses
		    (id, user_id, label, street, district, city, lat, lng, is_default)
		VALUES
		    (:id, :user_id, :label, :street, :district, :city, :lat, :lng, :is_default)
	`
	_, err := r.db.NamedExecContext(ctx, query, addr)
	if err != nil {
		return fmt.Errorf("AddressRepository.Create: %w", err)
	}
	return nil
}

func (r *AddressRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.UserAddress, error) {
	var addresses []model.UserAddress
	err := r.db.SelectContext(ctx, &addresses, `
		SELECT * FROM users.user_addresses
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("AddressRepository.FindByUserID: %w", err)
	}
	return addresses, nil
}

func (r *AddressRepository) FindByID(ctx context.Context, id, userID uuid.UUID) (*model.UserAddress, error) {
	var addr model.UserAddress
	err := r.db.GetContext(ctx, &addr, `
		SELECT * FROM users.user_addresses
		WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("AddressRepository.FindByID: %w", err)
	}
	return &addr, nil
}

func (r *AddressRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.GetContext(ctx, &count,
		`SELECT COUNT(*) FROM users.user_addresses WHERE user_id = $1`, userID)
	return count, err
}

func (r *AddressRepository) Update(ctx context.Context, addr *model.UserAddress) error {
	query := `
		UPDATE users.user_addresses
		SET label      = :label,
		    street     = :street,
		    district   = :district,
		    city       = :city,
		    lat        = :lat,
		    lng        = :lng,
		    updated_at = NOW()
		WHERE id = :id AND user_id = :user_id
	`
	_, err := r.db.NamedExecContext(ctx, query, addr)
	if err != nil {
		return fmt.Errorf("AddressRepository.Update: %w", err)
	}
	return nil
}

func (r *AddressRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM users.user_addresses WHERE id = $1 AND user_id = $2
	`, id, userID)
	return err
}

// UnsetAllDefault bỏ is_default của tất cả địa chỉ của user
func (r *AddressRepository) UnsetAllDefault(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users.user_addresses SET is_default = false WHERE user_id = $1
	`, userID)
	return err
}

// SetDefault đặt 1 địa chỉ là mặc định (unset tất cả trước)
func (r *AddressRepository) SetDefault(ctx context.Context, id, userID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`UPDATE users.user_addresses SET is_default = false WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE users.user_addresses
		SET is_default = true, updated_at = NOW()
		WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
