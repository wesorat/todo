package handler

import (
	"example/todo/internal/service"
	"log/slog"

	"github.com/gin-gonic/gin"
)

//  мидлвейр для идентификации
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
	api := router.Group("/api")
	api.Use(h.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllList)
			lists.PUT("/:id", h.updateList)
			lists.GET("/:id", h.getList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group("/:id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItem)
			}

		}
		items := api.Group("/items")
		{

			items.PUT("/:id", h.updateItem)
			items.GET("/:id", h.getItem)
			items.DELETE("/:id", h.deleteItem)
		}
	}

	return router

}
