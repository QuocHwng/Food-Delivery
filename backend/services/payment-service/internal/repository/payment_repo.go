package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"payment-service/internal/model"
)

type PaymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, p *model.Payment) error {
	query := `INSERT INTO payments.payments (id, order_id, amount, method, status)
		VALUES (:id, :order_id, :amount, :method, :status)`
	_, err := r.db.NamedExecContext(ctx, query, p)
	return err
}

func (r *PaymentRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) (*model.Payment, error) {
	var p model.Payment
	err := r.db.GetContext(ctx, &p, `SELECT * FROM payments.payments WHERE order_id = $1 ORDER BY created_at DESC LIMIT 1`, orderID)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	return &p, err
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, txnID *string) error {
	query := `UPDATE payments.payments SET status = $1, vnpay_txn_no = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, txnID, id)
	return err
}

func (r *PaymentRepository) Log(ctx context.Context, paymentID uuid.UUID, action, data string) error {
	logEntry := model.PaymentLog{
		ID:        uuid.New(),
		PaymentID: paymentID,
		Event:     action,
		Direction: "outbound", // defaulting to outbound for simplicity
		Payload:   fmt.Sprintf(`{"message": "%s"}`, data), // very simple json wrapper
	}
	query := `INSERT INTO payments.payment_logs (id, payment_id, event, direction, payload) VALUES (:id, :payment_id, :event, :direction, :payload)`
	_, err := r.db.NamedExecContext(ctx, query, logEntry)
	return err
}
