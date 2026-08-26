package handler

import (
	"github.com/gin-gonic/gin"
	"payment-service/internal/model"
	"payment-service/internal/service"
)

type PaymentHandler struct {
	svc *service.PaymentService
}

func NewPaymentHandler(svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

// CreatePayment is called by the frontend or order-service right after an order is placed
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req model.CreatePaymentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}
	req.IPAddress = c.ClientIP()

	url, err := h.svc.CreatePayment(c.Request.Context(), req)
	if err != nil {
		respondInternalError(c, err.Error())
		return
	}

	respondOK(c, gin.H{"checkout_url": url})
}

// VNPayIPN is the Server-to-Server webhook called by VNPay
func (h *PaymentHandler) VNPayIPN(c *gin.Context) {
	code, msg := h.svc.HandleIPN(c.Request.Context(), c.Request.URL.Query())
	c.JSON(200, gin.H{
		"RspCode": code,
		"Message": msg,
	})
}

// VNPayReturn is where the user is redirected after VNPay UI
func (h *PaymentHandler) VNPayReturn(c *gin.Context) {
	// For MVP, just redirect to frontend URL or show success JSON
	rspCode := c.Query("vnp_ResponseCode")
	if rspCode == "00" {
		c.JSON(200, gin.H{"success": true, "message": "Thanh toán thành công!"})
	} else {
		c.JSON(400, gin.H{"success": false, "message": "Thanh toán thất bại hoặc đã bị huỷ!"})
	}
}
