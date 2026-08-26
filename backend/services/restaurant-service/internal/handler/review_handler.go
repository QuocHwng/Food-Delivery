package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"restaurant-service/internal/model"
	"restaurant-service/internal/service"
)

type ReviewHandler struct {
	svc *service.ReviewService
}

func NewReviewHandler(svc *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{svc: svc}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	customerIDStr := c.GetString("user_id")
	customerID, _ := uuid.Parse(customerIDStr)
	
	orderIDStr := c.Param("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req model.CreateReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.CreateReview(c.Request.Context(), customerID, orderID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Review created successfully"})
}

func (h *ReviewHandler) GetRestaurantReviews(c *gin.Context) {
	restIDStr := c.Param("id")
	restID, err := uuid.Parse(restIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restaurant ID"})
		return
	}

	list, err := h.svc.GetRestaurantReviews(c.Request.Context(), restID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *ReviewHandler) ReplyReview(c *gin.Context) {
	ownerIDStr := c.GetString("user_id")
	ownerID, _ := uuid.Parse(ownerIDStr)

	reviewIDStr := c.Param("id")
	reviewID, err := uuid.Parse(reviewIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	var req model.ReplyReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.ReplyToReview(c.Request.Context(), ownerID, reviewID, req.Reply); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
