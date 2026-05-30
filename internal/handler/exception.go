package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) newErrorResponse(c *gin.Context, statusCode int, message string) {
	h.log.Error(message)
	c.AbortWithStatusJSON(statusCode, message)
}
