package handler

import (
	"example/todo/internal/service"
	"log/slog"

	"github.com/gin-gonic/gin"
)

// мидлвейр для идентификации
// ручки для списков и элементов

type Handler struct {
	service *service.Service
	log     *slog.Logger
}

func NewHandler(service *service.Service, log *slog.Logger) *Handler {
	return &Handler{service: service, log: log}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{

		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/logout", h.logout)
		auth.POST("/logout-all", h.logout_all)
		auth.POST("/refresh", h.refresh)
	}

	return router

}
