package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"api-gateway/internal/config"
	"api-gateway/internal/middleware"
	"api-gateway/internal/proxy"
)

func main() {
	cfg := config.Load()

	r := gin.Default()
	r.Use(middleware.SetupCORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "api-gateway",
			"routes": []string{
				"/api/auth", "/api/users",
				"/api/restaurants", "/api/menu-categories", "/api/menu-items",
				"/api/orders", "/api/coupons",
				"/api/payments",
				"/api/notifications", "/ws",
			},
		})
	})

	// Generate Handlers for each microservice
	userServiceProxy := proxy.ReverseProxy(cfg.UserServiceURL)
	restaurantServiceProxy := proxy.ReverseProxy(cfg.RestaurantServiceURL)
	orderServiceProxy := proxy.ReverseProxy(cfg.OrderServiceURL)
	paymentServiceProxy := proxy.ReverseProxy(cfg.PaymentServiceURL)
	notificationServiceProxy := proxy.ReverseProxy(cfg.NotificationServiceURL)

	// User Service Routes
	r.Any("/api/auth/*path", userServiceProxy)
	r.Any("/api/users/*path", userServiceProxy)

	// Restaurant Service Routes
	r.Any("/api/restaurants/*path", restaurantServiceProxy)
	r.Any("/api/menu-categories/*path", restaurantServiceProxy)
	r.Any("/api/menu-items/*path", restaurantServiceProxy)

	// Order Service Routes
	r.Any("/api/orders/*path", orderServiceProxy)
	r.Any("/api/coupons/*path", orderServiceProxy)

	// Payment Service Routes
	r.Any("/api/payments/*path", paymentServiceProxy)

	// Notification Service Routes
	r.Any("/api/notifications/*path", notificationServiceProxy)
	r.Any("/ws", notificationServiceProxy) // WebSocket support is handled transparently by httputil.ReverseProxy in Go 1.12+

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 API Gateway running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("❌ Gateway error: %v", err)
	}
}
