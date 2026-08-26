package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"restaurant-service/internal/model"
	"restaurant-service/internal/service"
)

type RestaurantHandler struct {
	svc *service.RestaurantService
}

func NewRestaurantHandler(svc *service.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{svc: svc}
}

// ─── Public Endpoints ────────────────────────────────────────────────────────

func (h *RestaurantHandler) ListRestaurants(c *gin.Context) {
	list, err := h.svc.ListRestaurants(c.Request.Context())
	if err != nil {
		respondInternalError(c)
		return
	}
	respondOK(c, list)
}

func (h *RestaurantHandler) GetRestaurantDetail(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondBadRequest(c, "invalid restaurant id")
		return
	}

	rest, err := h.svc.GetRestaurantDetail(c.Request.Context(), id)
	if err != nil {
		respondNotFound(c, "restaurant not found")
		return
	}
	respondOK(c, rest)
}

func (h *RestaurantHandler) GetMenu(c *gin.Context) {
	restID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondBadRequest(c, "invalid restaurant id")
		return
	}
	menu, err := h.svc.GetFullMenu(c.Request.Context(), restID)
	if err != nil {
		respondInternalError(c)
		return
	}
	respondOK(c, menu)
}

// ─── Owner Endpoints ─────────────────────────────────────────────────────────

func (h *RestaurantHandler) CreateRestaurant(c *gin.Context) {
	ownerID, _ := uuid.Parse(c.GetString("user_id"))
	role := c.GetString("user_role")

	var req model.CreateRestaurantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	rest, err := h.svc.CreateRestaurant(c.Request.Context(), ownerID, role, req)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			respondForbidden(c, "only restaurant_owner can create")
			return
		}
		respondInternalError(c)
		return
	}
	respondCreated(c, rest)
}

func (h *RestaurantHandler) CreateMenuCategory(c *gin.Context) {
	ownerID, _ := uuid.Parse(c.GetString("user_id"))
	restID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondBadRequest(c, "invalid restaurant id")
		return
	}

	var req model.CreateMenuCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	cat, err := h.svc.CreateMenuCategory(c.Request.Context(), ownerID, restID, req)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			respondForbidden(c, "you do not own this restaurant")
			return
		}
		respondInternalError(c)
		return
	}
	respondCreated(c, cat)
}

func (h *RestaurantHandler) CreateMenuItem(c *gin.Context) {
	ownerID, _ := uuid.Parse(c.GetString("user_id"))
	restID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondBadRequest(c, "invalid restaurant id")
		return
	}

	var req model.CreateMenuItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	item, err := h.svc.CreateMenuItem(c.Request.Context(), ownerID, restID, req)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			respondForbidden(c, "you do not own this restaurant")
			return
		}
		respondInternalError(c)
		return
	}
	respondCreated(c, item)
}
