package handler

import "github.com/gin-gonic/gin"

func SetupRouter(h *RestaurantHandler, reviewHandler *ReviewHandler, jwtSecret string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "restaurant-service"})
	})

	api := r.Group("/api/restaurants")
	{
		api.GET("", h.ListRestaurants)
		api.GET("/:id", h.GetRestaurantDetail)
		api.GET("/:id/menu", h.GetMenu)
		api.GET("/:id/ratings", reviewHandler.GetRestaurantReviews)

		protected := api.Group("")
		protected.Use(JWTMiddleware(jwtSecret))
		{
			protected.POST("", h.CreateRestaurant)

			// Review related
			protected.POST("/orders/:order_id/ratings", reviewHandler.CreateReview)
			protected.POST("/ratings/:id/reply", reviewHandler.ReplyReview)
			protected.PUT("/ratings/:id/reply", reviewHandler.ReplyReview)

			// Only owners (or admins)
			owners := protected.Group("/")
			owners.Use(RequireRole("restaurant_owner", "admin"))
			{
				owners.POST("/:id/menu-categories", h.CreateMenuCategory)
				owners.POST("/:id/menu-items", h.CreateMenuItem)
			}
		}
	}

	return r
}
