package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"order-service/internal/model"
	"order-service/internal/service"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) PlaceOrder(c *gin.Context) {
	customerID, _ := uuid.Parse(c.GetString("user_id"))

	var req model.CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	order, err := h.svc.PlaceOrder(c.Request.Context(), customerID, req)
	if err != nil {
		respondInternalError(c, err.Error())
		return
	}
	respondCreated(c, order)
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))
	role := c.GetString("user_role")

	var list []model.Order
	var err error

	if role == "customer" {
		list, err = h.svc.GetCustomerOrders(c.Request.Context(), userID)
	} else if role == "restaurant_owner" {
		// NOTE: In a real scenario, we should get the restaurant_id owned by this user
		// For MVP, passing restaurant_id in query if role=owner
		restIDStr := c.Query("restaurant_id")
		if restIDStr == "" {
			respondBadRequest(c, "restaurant_id query param required for owners")
			return
		}
		restID, _ := uuid.Parse(restIDStr)
		list, err = h.svc.GetRestaurantOrders(c.Request.Context(), restID)
	}

	if err != nil {
		respondInternalError(c, err.Error())
		return
	}
	respondOK(c, list)
}

func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondBadRequest(c, "invalid order id")
		return
	}
	userID, _ := uuid.Parse(c.GetString("user_id"))
	role := c.GetString("user_role")

	order, err := h.svc.GetOrderDetail(c.Request.Context(), orderID, userID, role)
	if err != nil {
		respondNotFound(c, "order not found or unauthorized")
		return
	}
	respondOK(c, order)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	orderID, _ := uuid.Parse(c.Param("id"))
	userID, _ := uuid.Parse(c.GetString("user_id"))

	var req model.UpdateStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateStatus(c.Request.Context(), orderID, userID, req.Status, req.Note); err != nil {
		respondInternalError(c, err.Error())
		return
	}
	respondOK(c, gin.H{"message": "status updated successfully"})
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	orderID, _ := uuid.Parse(c.Param("id"))
	userID, _ := uuid.Parse(c.GetString("user_id"))

	var req model.CancelOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	if err := h.svc.CancelOrder(c.Request.Context(), orderID, userID, req.Reason); err != nil {
		respondBadRequest(c, err.Error())
		return
	}
	respondOK(c, gin.H{"message": "order cancelled successfully"})
}
