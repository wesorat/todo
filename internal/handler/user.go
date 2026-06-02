package handler

import (
	"errors"
	"example/todo/internal/domain"
	"example/todo/internal/repository"
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
		if errors.Is(err, repository.ErrNameExists) {
			h.newErrorResponse(c, http.StatusConflict, "name already in use")
			return
		}
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]int{
		"id": id,
	})

}

type signInInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input domain.CreateUser
	if err := c.BindJSON(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.service.User.GetUser(input.Name, input.Password)
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	_ = user

	c.JSON(http.StatusOK, map[string]string{
		"message": "logged is succesfully",
	})
}
