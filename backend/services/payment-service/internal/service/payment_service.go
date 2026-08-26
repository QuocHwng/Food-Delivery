package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"payment-service/internal/config"
	"payment-service/internal/model"
	"payment-service/internal/repository"
)

type PaymentService struct {
	cfg  *config.Config
	repo *repository.PaymentRepository
}

func NewPaymentService(cfg *config.Config, repo *repository.PaymentRepository) *PaymentService {
	return &PaymentService{cfg: cfg, repo: repo}
}

// CreatePayment generates a payment record and returns the VNPay checkout URL if method is 'vnpay'
func (s *PaymentService) CreatePayment(ctx context.Context, req model.CreatePaymentReq) (string, error) {
	// Check if already exists to prevent duplicate payment intents
	existing, err := s.repo.FindByOrderID(ctx, req.OrderID)
	if err == nil && existing != nil && existing.Status == "success" {
		return "", fmt.Errorf("order already paid")
	}

	payment := &model.Payment{
		ID:         uuid.New(),
		OrderID:    req.OrderID,
		Amount:     req.Amount,
		Method:     req.Method,
		Status:     "pending",
	}

	if err := s.repo.CreatePayment(ctx, payment); err != nil {
		return "", err
	}

	s.repo.Log(ctx, payment.ID, "CREATED", fmt.Sprintf("Method: %s", req.Method))

	if req.Method == "cod" {
		// Cash on delivery doesn't require a checkout URL
		return "COD_CONFIRMED", nil
	}

	// Generate VNPay URL
	checkoutUrl := GenerateVNPayURL(s.cfg, payment, req.IPAddress)
	s.repo.Log(ctx, payment.ID, "VNPAY_URL_GENERATED", checkoutUrl)
	return checkoutUrl, nil
}

// HandleIPN processes the server-to-server callback from VNPay
func (s *PaymentService) HandleIPN(ctx context.Context, query url.Values) (string, string) {
	// 1. Verify signature
	if !VerifyVNPaySignature(s.cfg.VnpHashSecret, query) {
		return "97", "Invalid signature"
	}

	paymentIDStr := query.Get("vnp_TxnRef")
	paymentID, err := uuid.Parse(paymentIDStr)
	if err != nil {
		return "01", "Order not found"
	}

	rspCode := query.Get("vnp_ResponseCode")
	txnNo := query.Get("vnp_TransactionNo")

	// Dump query to JSON for logging
	queryData, _ := json.Marshal(query)
	s.repo.Log(ctx, paymentID, "VNPAY_IPN_RECEIVED", string(queryData))

	status := "failed"
	if rspCode == "00" {
		status = "success"
	}

	if err := s.repo.UpdateStatus(ctx, paymentID, status, &txnNo); err != nil {
		return "99", "Unknown error"
	}

	// TODO: Publish RabbitMQ message `payment.success` so Order Service updates order status

	return "00", "Confirm Success"
}
