package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"notification-service/internal/model"
)

type NotificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(ctx context.Context, n *model.Notification) error {
	query := `INSERT INTO notifications.notifications (id, user_id, title, message, type, data)
		VALUES (:id, :user_id, :title, :message, :type, :data)`
	_, err := r.db.NamedExecContext(ctx, query, n)
	return err
}

func (r *NotificationRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]model.Notification, error) {
	var list []model.Notification
	query := `SELECT * FROM notifications.notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT 50`
	err := r.db.SelectContext(ctx, &list, query, userID)
	return list, err
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `UPDATE notifications.notifications SET is_read = true WHERE id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, id, userID)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE notifications.notifications SET is_read = true WHERE user_id = $1 AND is_read = false`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
