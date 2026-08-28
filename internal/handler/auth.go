package handler

import (
	"errors"
	"net/http"

	"github.com/wesorat/todo/internal/domain"
	"github.com/wesorat/todo/internal/repository"

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
	tokens, err := h.service.Auth.SignIn(c.Request.Context(), input.Name, input.Password)
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
		"access_token": tokens.JWT,
	})
}

func (h *Handler) refresh(c *gin.Context) {
	refresh_token, err := c.Cookie("refresh_token")
	if err != nil {
		h.newErrorResponse(c, http.StatusUnauthorized, "refresh_token is missing")
		return
	}
	access_token, err := h.service.Auth.RenewalJWT(c.Request.Context(), refresh_token)
	if err != nil {
		h.newErrorResponse(c, http.StatusUnauthorized, "refresh_token is invalid")
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"token": access_token,
	})
}

func (h *Handler) logout(c *gin.Context) {
	refresh_token, err := c.Cookie("refresh_token")
	if err != nil {
		h.newErrorResponse(c, http.StatusUnauthorized, "refresh_token is missing")
		return
	}
	if err := h.service.Auth.Logout(c.Request.Context(), refresh_token); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "couldnt remove refresh_token")
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})
	c.JSON(http.StatusOK, map[string]string{
		"message": "logout successfully",
	})
}

func (h *Handler) logout_all(c *gin.Context) {
	refresh_token, err := c.Cookie("refresh_token")
	if err != nil {
		h.newErrorResponse(c, http.StatusUnauthorized, "refresh_token is missing")
		return
	}
	if err := h.service.Auth.LogoutAll(c.Request.Context(), refresh_token); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "couldnt remove refresh_token")
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})
	c.JSON(http.StatusOK, map[string]string{
		"message": "logout all successfully",
	})
}
