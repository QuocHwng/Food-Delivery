package handler

import "github.com/gin-gonic/gin"

func SetupRouter(h *NotificationHandler, jwtSecret string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "notification-service"})
	})

	// WebSockets
	r.GET("/ws", AuthMiddleware(jwtSecret), h.ServeWS)

	// REST APIs
	api := r.Group("/api/notifications")
	api.Use(AuthMiddleware(jwtSecret))
	{
		api.GET("/my", h.GetMyNotifications)
		api.PATCH("/:id/read", h.MarkAsRead)
		api.PATCH("/read-all", h.MarkAllAsRead)
	}

	return r
}
