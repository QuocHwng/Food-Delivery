package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"order-service/internal/model"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Execute in a transaction since it inserts into orders, order_items, and order_status_history
func (r *OrderRepository) CreateOrder(ctx context.Context, order *model.Order, items []model.OrderItem, history *model.OrderStatusHistory) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // ignored if tx.Commit() is called first

	// Insert Order
	oq := `INSERT INTO orders.orders 
		(id, customer_id, restaurant_id, address_id, status, subtotal, discount_amount, delivery_fee, total, note)
		VALUES (:id, :customer_id, :restaurant_id, :address_id, :status, :subtotal, :discount_amount, :delivery_fee, :total, :note)`
	if _, err := tx.NamedExecContext(ctx, oq, order); err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	// Insert Items
	iq := `INSERT INTO orders.order_items 
		(id, order_id, menu_item_id, name_snapshot, price_snapshot, quantity)
		VALUES (:id, :order_id, :menu_item_id, :name_snapshot, :price_snapshot, :quantity)`
	for _, it := range items {
		if _, err := tx.NamedExecContext(ctx, iq, it); err != nil {
			return fmt.Errorf("insert order item: %w", err)
		}
	}

	// Insert History
	hq := `INSERT INTO orders.order_status_history (id, order_id, to_status, changed_by)
		VALUES (:id, :order_id, :to_status, :changed_by)`
	if _, err := tx.NamedExecContext(ctx, hq, history); err != nil {
		return fmt.Errorf("insert history: %w", err)
	}

	return tx.Commit()
}

func (r *OrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	var o model.Order
	err := r.db.GetContext(ctx, &o, `SELECT * FROM orders.orders WHERE id = $1`, id)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}

	// Fetch items
	err = r.db.SelectContext(ctx, &o.Items, `SELECT * FROM orders.order_items WHERE order_id = $1`, id)
	return &o, err
}

func (r *OrderRepository) FindByCustomer(ctx context.Context, customerID uuid.UUID) ([]model.Order, error) {
	var list []model.Order
	err := r.db.SelectContext(ctx, &list, `SELECT * FROM orders.orders WHERE customer_id = $1 ORDER BY created_at DESC`, customerID)
	return list, err
}

func (r *OrderRepository) FindByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]model.Order, error) {
	var list []model.Order
	err := r.db.SelectContext(ctx, &list, `SELECT * FROM orders.orders WHERE restaurant_id = $1 ORDER BY created_at DESC`, restaurantID)
	return list, err
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, h *model.OrderStatusHistory) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update order
	if _, err := tx.ExecContext(ctx, `UPDATE orders.orders SET status = $1, updated_at = NOW() WHERE id = $2`, status, id); err != nil {
		return err
	}

	// Insert history
	hq := `INSERT INTO orders.order_status_history (id, order_id, from_status, to_status, changed_by, note)
		VALUES (:id, :order_id, :from_status, :to_status, :changed_by, :note)`
	if _, err := tx.NamedExecContext(ctx, hq, h); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *OrderRepository) CancelOrder(ctx context.Context, id uuid.UUID, cancelledBy uuid.UUID, reason string, h *model.OrderStatusHistory) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update order
	uq := `UPDATE orders.orders SET status = 'cancelled', cancelled_by = $1, cancel_reason = $2, updated_at = NOW() WHERE id = $3`
	if _, err := tx.ExecContext(ctx, uq, cancelledBy, reason, id); err != nil {
		return err
	}

	// Insert history
	hq := `INSERT INTO orders.order_status_history (id, order_id, from_status, to_status, changed_by, note)
		VALUES (:id, :order_id, :from_status, :to_status, :changed_by, :note)`
	if _, err := tx.NamedExecContext(ctx, hq, h); err != nil {
		return err
	}

	return tx.Commit()
}
