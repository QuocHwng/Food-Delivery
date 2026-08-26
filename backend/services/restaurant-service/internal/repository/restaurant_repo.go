package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"restaurant-service/internal/model"
)

type RestaurantRepository struct {
	db *sqlx.DB
}

func NewRestaurantRepository(db *sqlx.DB) *RestaurantRepository {
	return &RestaurantRepository{db: db}
}

func (r *RestaurantRepository) Create(ctx context.Context, rest *model.Restaurant) error {
	query := `
		INSERT INTO restaurants.restaurants 
		(id, owner_id, name, description, address, phone, min_order_value, delivery_fee)
		VALUES (:id, :owner_id, :name, :description, :address, :phone, :min_order_value, :delivery_fee)
	`
	_, err := r.db.NamedExecContext(ctx, query, rest)
	if err != nil {
		return fmt.Errorf("RestaurantRepository.Create: %w", err)
	}
	return nil
}

func (r *RestaurantRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Restaurant, error) {
	var rest model.Restaurant
	err := r.db.GetContext(ctx, &rest, `SELECT * FROM restaurants.restaurants WHERE id = $1 AND is_active = true`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &rest, nil
}

func (r *RestaurantRepository) FindByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]model.Restaurant, error) {
	var list []model.Restaurant
	err := r.db.SelectContext(ctx, &list, `SELECT * FROM restaurants.restaurants WHERE owner_id = $1`, ownerID)
	return list, err
}

func (r *RestaurantRepository) ListAll(ctx context.Context) ([]model.Restaurant, error) {
	var list []model.Restaurant
	err := r.db.SelectContext(ctx, &list, `SELECT * FROM restaurants.restaurants WHERE is_active = true AND is_open = true ORDER BY avg_rating DESC`)
	return list, err
}
