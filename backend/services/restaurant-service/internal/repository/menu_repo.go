package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"restaurant-service/internal/model"
)

type MenuRepository struct {
	db *sqlx.DB
}

func NewMenuRepository(db *sqlx.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

func (r *MenuRepository) CreateCategory(ctx context.Context, cat *model.MenuCategory) error {
	query := `
		INSERT INTO restaurants.menu_categories (id, restaurant_id, name, display_order)
		VALUES (:id, :restaurant_id, :name, :display_order)
	`
	_, err := r.db.NamedExecContext(ctx, query, cat)
	return err
}

func (r *MenuRepository) CreateItem(ctx context.Context, item *model.MenuItem) error {
	query := `
		INSERT INTO restaurants.menu_items (id, restaurant_id, category_id, name, description, price, image_url)
		VALUES (:id, :restaurant_id, :category_id, :name, :description, :price, :image_url)
	`
	_, err := r.db.NamedExecContext(ctx, query, item)
	return err
}

func (r *MenuRepository) GetCategoriesByRestaurant(ctx context.Context, restID uuid.UUID) ([]model.MenuCategory, error) {
	var list []model.MenuCategory
	err := r.db.SelectContext(ctx, &list, `SELECT * FROM restaurants.menu_categories WHERE restaurant_id = $1 AND is_active = true ORDER BY display_order`, restID)
	return list, err
}

func (r *MenuRepository) GetItemsByRestaurant(ctx context.Context, restID uuid.UUID) ([]model.MenuItem, error) {
	var list []model.MenuItem
	err := r.db.SelectContext(ctx, &list, `SELECT * FROM restaurants.menu_items WHERE restaurant_id = $1 AND is_available = true ORDER BY display_order`, restID)
	return list, err
}
