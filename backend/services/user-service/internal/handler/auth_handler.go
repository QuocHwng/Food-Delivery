package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"user-service/internal/model"
	"user-service/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// Register godoc
// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	resp, err := h.authSvc.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			respondConflict(c, "email already registered")
			return
		}
		respondInternalError(c)
		return
	}

	respondCreated(c, resp)
}

// Login godoc
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	resp, err := h.authSvc.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCreds) {
			respondUnauthorized(c, "invalid email or password")
			return
		}
		respondInternalError(c)
		return
	}

	respondOK(c, resp)
}

// Refresh godoc
// POST /api/auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req model.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	resp, err := h.authSvc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		respondUnauthorized(c, "invalid or expired refresh token")
		return
	}

	respondOK(c, resp)
}

// Logout godoc
// POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var req model.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	if err := h.authSvc.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		respondInternalError(c)
		return
	}

	respondOK(c, gin.H{"message": "logged out successfully"})
}
