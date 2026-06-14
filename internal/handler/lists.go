package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(c *gin.Context) {
	val, ok := c.Get("user_id")
	if !ok {
		h.newErrorResponse(c, http.StatusUnauthorized, "")
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"user_id": val,
	})
}
