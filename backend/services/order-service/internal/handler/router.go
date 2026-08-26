package handler

import "github.com/gin-gonic/gin"

func SetupRouter(h *OrderHandler, dashHandler *DashboardHandler, jwtSecret string) *gin.Engine {
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

		// Dashboard & Stats routes (for Restaurant Owner)
		dash := protected.Group("/orders/restaurant/:id")
		dash.Use(RequireRole("restaurant_owner", "admin"))
		{
			dash.GET("/dashboard", dashHandler.GetOverview)
			dash.GET("/active", dashHandler.GetActiveOrders)
			dash.GET("/stats/revenue", dashHandler.GetRevenueStats)
			dash.GET("/stats/top-items", dashHandler.GetTopItems)
			dash.GET("/stats/orders-count", dashHandler.GetOrderCounts)
		}
	}

	return r
}
