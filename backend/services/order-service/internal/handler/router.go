package handler

import "github.com/gin-gonic/gin"

func SetupRouter(h *OrderHandler, jwtSecret string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "order-service"})
	})

	api := r.Group("/api")
	
	// Protected APIs
	protected := api.Group("/")
	protected.Use(JWTMiddleware(jwtSecret))
	{
		// Customer routes
		protected.POST("/orders", RequireRole("customer"), h.PlaceOrder)
		
		// Shared routes (Customer / Owner)
		protected.GET("/orders", h.GetMyOrders)
		protected.GET("/orders/:id", h.GetOrderDetail)
		protected.PATCH("/orders/:id/cancel", h.CancelOrder)

		// Owner routes
		protected.PATCH("/orders/:id/status", RequireRole("restaurant_owner", "admin"), h.UpdateStatus)
	}

	return r
}
