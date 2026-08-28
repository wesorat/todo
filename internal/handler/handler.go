package handler

import (
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wesorat/todo/internal/service"

	"github.com/gin-gonic/gin"
)

//  мидлвейр для идентификации
// ручки для списков и элементов

type Handler struct {
	service *service.Service
	log     *slog.Logger
	redis   *redis.Client
}

func NewHandler(service *service.Service, log *slog.Logger, redis *redis.Client) *Handler {
	return &Handler{service: service, log: log, redis: redis}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		authLimiter := newRedisRateLimiter(h.redis, 5, time.Minute)
		auth.POST("/sign-up", rateLimit(authLimiter), h.signUp)
		auth.POST("/sign-in", rateLimit(authLimiter), h.signIn)
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
