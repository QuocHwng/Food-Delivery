package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"user-service/internal/model"
	"user-service/internal/service"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// ─── Profile ──────────────────────────────────────────────────────────────────

// GetProfile godoc
// GET /api/users/me
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := mustParseUUID(c, "user_id")

	user, err := h.userSvc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		respondNotFound(c, "user not found")
		return
	}
	respondOK(c, user)
}

// UpdateProfile godoc
// PUT /api/users/me
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := mustParseUUID(c, "user_id")

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	user, err := h.userSvc.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		respondInternalError(c)
		return
	}
	respondOK(c, user)
}

// ─── Addresses ────────────────────────────────────────────────────────────────

// GetAddresses godoc
// GET /api/users/me/addresses
func (h *UserHandler) GetAddresses(c *gin.Context) {
	userID := mustParseUUID(c, "user_id")

	addresses, err := h.userSvc.GetAddresses(c.Request.Context(), userID)
	if err != nil {
		respondInternalError(c)
		return
	}
	respondOK(c, addresses)
}

// CreateAddress godoc
// POST /api/users/me/addresses
func (h *UserHandler) CreateAddress(c *gin.Context) {
	userID := mustParseUUID(c, "user_id")

	var req model.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	addr, err := h.userSvc.CreateAddress(c.Request.Context(), userID, req)
	if err != nil {
		respondInternalError(c)
		return
	}
	respondCreated(c, addr)
}

// UpdateAddress godoc
// PUT /api/users/me/addresses/:id
func (h *UserHandler) UpdateAddress(c *gin.Context) {
	userID := mustParseUUID(c, "user_id")
	addrID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondBadRequest(c, "invalid address id")
		return
	}

	var req model.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	addr, err := h.userSvc.UpdateAddress(c.Request.Context(), userID, addrID, req)
	if err != nil {
		respondInternalError(c)
		return
	}
	respondOK(c, addr)
}

// DeleteAddress godoc
// DELETE /api/users/me/addresses/:id
func (h *UserHandler) DeleteAddress(c *gin.Context) {
	userID := mustParseUUID(c, "user_id")
	addrID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondBadRequest(c, "invalid address id")
		return
	}

	if err := h.userSvc.DeleteAddress(c.Request.Context(), userID, addrID); err != nil {
		respondInternalError(c)
		return
	}
	respondOK(c, gin.H{"message": "address deleted"})
}

// SetDefaultAddress godoc
// PATCH /api/users/me/addresses/:id/default
func (h *UserHandler) SetDefaultAddress(c *gin.Context) {
	userID := mustParseUUID(c, "user_id")
	addrID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondBadRequest(c, "invalid address id")
		return
	}

	if err := h.userSvc.SetDefaultAddress(c.Request.Context(), userID, addrID); err != nil {
		respondInternalError(c)
		return
	}
	respondOK(c, gin.H{"message": "default address updated"})
}

// ─── Helper ───────────────────────────────────────────────────────────────────

func mustParseUUID(c *gin.Context, key string) uuid.UUID {
	id, _ := uuid.Parse(c.GetString(key))
	return id
}
