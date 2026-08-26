package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"user-service/internal/service"
)

// JWTMiddleware xác thực Bearer token, inject user_id và user_role vào context
func JWTMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			respondUnauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := service.ValidateAccessToken(tokenStr, jwtSecret)
		if err != nil {
			respondUnauthorized(c, "invalid or expired access token")
			c.Abort()
			return
		}

		// Inject vào context để handler dùng
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}
