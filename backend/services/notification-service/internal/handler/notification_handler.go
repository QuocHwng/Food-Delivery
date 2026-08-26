package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"notification-service/internal/service"
	"notification-service/internal/websocket"
)

type NotificationHandler struct {
	svc *service.NotificationService
	hub *websocket.Hub
}

func NewNotificationHandler(svc *service.NotificationService, hub *websocket.Hub) *NotificationHandler {
	return &NotificationHandler{svc: svc, hub: hub}
}

func (h *NotificationHandler) GetMyNotifications(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	list, err := h.svc.GetMyNotifications(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	notifIDStr := c.Param("id")
	notifID, err := uuid.Parse(notifIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	if err := h.svc.MarkAsRead(c.Request.Context(), notifID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	if err := h.svc.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *NotificationHandler) ServeWS(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	websocket.ServeWs(h.hub, c.Writer, c.Request, userIDStr)
}
