package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) userIdentity(c *gin.Context) {
	access := c.GetHeader("Authorization")
	if access == "" {
		h.newErrorResponse(c, http.StatusUnauthorized, "missing access token in header")
		return
	}
	user_id, err := h.service.Auth.ParseJWT(access)
	if err != nil {
		h.newErrorResponse(c, http.StatusUnauthorized, "access token is invalid")
		return
	}
	c.Set("user_id", user_id)
}
