package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"user-service/internal/model"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users.users (id, name, email, phone, password_hash, role)
		VALUES (:id, :name, :email, :phone, :password_hash, :role)
	`
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("UserRepository.Create: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user,
		`SELECT * FROM users.users WHERE email = $1 AND is_active = true`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("UserRepository.FindByEmail: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.GetContext(ctx, &user,
		`SELECT * FROM users.users WHERE id = $1 AND is_active = true`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("UserRepository.FindByID: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users.users
		SET name       = :name,
		    phone      = :phone,
		    avatar_url = :avatar_url,
		    updated_at = NOW()
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("UserRepository.Update: %w", err)
	}
	return nil
}
