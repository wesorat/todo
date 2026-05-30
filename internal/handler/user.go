package handler

import (
	"example/todo/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var user domain.CreateUser
	if err := c.BindJSON(&user); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.service.User.CreateUser(user)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

// TODO сделать уникальные name
type signInInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *Handler) signIn(c *gin.Context) {

}
