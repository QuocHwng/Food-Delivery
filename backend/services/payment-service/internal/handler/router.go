package handler

import "github.com/gin-gonic/gin"

func SetupRouter(h *PaymentHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "payment-service"})
	})

	api := r.Group("/api/payments")
	{
		api.POST("/create", h.CreatePayment)
		
		// VNPay Webhooks
		api.GET("/vnpay/ipn", h.VNPayIPN)
		api.GET("/vnpay/return", h.VNPayReturn)
	}

	return r
}
