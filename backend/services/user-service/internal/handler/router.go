package handler

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter đăng ký tất cả routes cho User Service
func SetupRouter(
	authHandler *AuthHandler,
	userHandler *UserHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// ── Health check ─────────────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "user-service"})
	})

	api := r.Group("/api")

	// ── Auth — public ─────────────────────────────────────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/logout", authHandler.Logout)
	}

	// ── Users — protected ─────────────────────────────────────────────────────
	users := api.Group("/users")
	users.Use(JWTMiddleware(jwtSecret))
	{
		users.GET("/me", userHandler.GetProfile)
		users.PUT("/me", userHandler.UpdateProfile)

		addresses := users.Group("/me/addresses")
		{
			addresses.GET("", userHandler.GetAddresses)
			addresses.POST("", userHandler.CreateAddress)
			addresses.PUT("/:id", userHandler.UpdateAddress)
			addresses.DELETE("/:id", userHandler.DeleteAddress)
			addresses.PATCH("/:id/default", userHandler.SetDefaultAddress)
		}
	}

	return r
}
