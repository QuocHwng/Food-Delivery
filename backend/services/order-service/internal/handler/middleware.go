package handler

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func JWTMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			respondUnauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		
		token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			respondUnauthorized(c, "invalid or expired access token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			respondUnauthorized(c, "invalid token claims")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("user_role")
		for _, r := range allowedRoles {
			if r == role {
				c.Next()
				return
			}
		}
		respondForbidden(c, "you do not have permission")
		c.Abort()
	}
}
