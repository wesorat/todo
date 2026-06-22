package handler

import (
	"net/http"
	"strconv"

	"github.com/wesorat/todo/internal/domain"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createItem(c *gin.Context) {
	user_id := c.GetInt("user_id")
	if user_id == 0 {
		h.newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var item domain.CreateItem

	list_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := c.Bind(&item); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	item.ListID = list_id
	item_id, err := h.service.Item.Create(user_id, item)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"item_id": item_id,
	})
}
func (h *Handler) getItem(c *gin.Context) {
	user_id := c.GetInt("user_id")
	if user_id == 0 {
		h.newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	item_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	item, err := h.service.Item.Get(user_id, item_id)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"item": item,
	})
}
func (h *Handler) getAllItem(c *gin.Context) {
	user_id := c.GetInt("user_id")
	if user_id == 0 {
		h.newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	list_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	items, err := h.service.Item.GetAll(user_id, list_id)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"items": items,
	})
}

func (h *Handler) updateItem(c *gin.Context) {
	user_id := c.GetInt("user_id")
	if user_id == 0 {
		h.newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var item domain.UpdateItem

	item_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := c.Bind(&item); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Item.Update(user_id, item_id, item.Title, item.Description, item.Done); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"status": "updated",
	})
}

func (h *Handler) deleteItem(c *gin.Context) {
	user_id := c.GetInt("user_id")
	if user_id == 0 {
		h.newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	item_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := h.service.Item.Delete(user_id, item_id); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"status": "deleted",
	})
}
