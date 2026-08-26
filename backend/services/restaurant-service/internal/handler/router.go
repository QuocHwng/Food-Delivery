package handler

import "github.com/gin-gonic/gin"

func SetupRouter(h *RestaurantHandler, jwtSecret string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "restaurant-service"})
	})

	api := r.Group("/api")

	// Public APIs
	api.GET("/restaurants", h.ListRestaurants)
	api.GET("/restaurants/:id", h.GetRestaurantDetail)
	api.GET("/restaurants/:id/menu", h.GetMenu)

	// Protected APIs
	protected := api.Group("/")
	protected.Use(JWTMiddleware(jwtSecret))
	{
		// Only owners (or admins)
		owners := protected.Group("/")
		owners.Use(RequireRole("restaurant_owner", "admin"))
		{
			owners.POST("/restaurants", h.CreateRestaurant)
			owners.POST("/restaurants/:id/menu-categories", h.CreateMenuCategory)
			owners.POST("/restaurants/:id/menu-items", h.CreateMenuItem)
		}
	}

	return r
}
