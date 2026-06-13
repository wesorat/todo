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
	id, err := h.service.Auth.CreateUser(user)
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
	tokens, err := h.service.Auth.SignIn(input.Name, input.Password)
	if err != nil {
		h.newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		MaxAge:   30 * 24 * 60 * 60,
		HttpOnly: true,
	})

	c.JSON(http.StatusOK, map[string]string{
		"access": tokens.JWT,
	})
}

// обновление токена
func (h *Handler) refresh(c *gin.Context) {}

// отзыв токена, очистка куки и удаление с бд
func (h *Handler) logout(c *gin.Context) {}

func (h *Handler) logout_all(c *gin.Context) {}
