package handler

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func respondOK(c *gin.Context, data interface{}) {
	c.JSON(200, APIResponse{Success: true, Data: data})
}

func respondCreated(c *gin.Context, data interface{}) {
	c.JSON(201, APIResponse{Success: true, Data: data})
}

func respondBadRequest(c *gin.Context, msg string) {
	c.JSON(400, APIResponse{Success: false, Message: msg})
}

func respondUnauthorized(c *gin.Context, msg string) {
	c.JSON(401, APIResponse{Success: false, Message: msg})
}

func respondForbidden(c *gin.Context, msg string) {
	c.JSON(403, APIResponse{Success: false, Message: msg})
}

func respondNotFound(c *gin.Context, msg string) {
	c.JSON(404, APIResponse{Success: false, Message: msg})
}

func respondInternalError(c *gin.Context, msg string) {
	c.JSON(500, APIResponse{Success: false, Message: msg})
}
