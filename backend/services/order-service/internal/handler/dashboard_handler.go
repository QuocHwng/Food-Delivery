package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"order-service/internal/service"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) GetOverview(c *gin.Context) {
	restID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
		return
	}

	data, err := h.svc.GetOverviewToday(c.Request.Context(), restID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *DashboardHandler) GetActiveOrders(c *gin.Context) {
	restID, _ := uuid.Parse(c.Param("id"))
	data, err := h.svc.GetActiveOrders(c.Request.Context(), restID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *DashboardHandler) GetRevenueStats(c *gin.Context) {
	restID, _ := uuid.Parse(c.Param("id"))
	
	fromStr := c.Query("from")
	toStr := c.Query("to")
	groupBy := c.Query("group_by")

	// Default to last 30 days if not provided
	now := time.Now()
	from := now.AddDate(0, 0, -30)
	to := now

	if fromStr != "" {
		if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = t
		}
	}
	if toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = t
		}
	}

	data, err := h.svc.GetRevenueStats(c.Request.Context(), restID, from, to, groupBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *DashboardHandler) GetTopItems(c *gin.Context) {
	restID, _ := uuid.Parse(c.Param("id"))
	data, err := h.svc.GetTopItems(c.Request.Context(), restID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *DashboardHandler) GetOrderCounts(c *gin.Context) {
	restID, _ := uuid.Parse(c.Param("id"))
	data, err := h.svc.GetOrderCounts(c.Request.Context(), restID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
