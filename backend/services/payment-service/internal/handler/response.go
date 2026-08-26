package handler

import "github.com/gin-gonic/gin"

func respondOK(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"success": true, "data": data})
}

func respondBadRequest(c *gin.Context, msg string) {
	c.JSON(400, gin.H{"success": false, "message": msg})
}

func respondInternalError(c *gin.Context, msg string) {
	c.JSON(500, gin.H{"success": false, "message": msg})
}
