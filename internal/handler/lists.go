package handler

import (
	"net/http"
	"strconv"

	"github.com/wesorat/todo/internal/domain"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(c *gin.Context) {
	user_id := c.GetInt("user_id")
	if user_id == 0 {
		h.newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	var input domain.CreateList

	if err := c.Bind(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	input.UserID = user_id
	list_id, err := h.service.List.Create(input)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"list_id": list_id,
	})

}

func (h *Handler) getList(c *gin.Context) {
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
	list, err := h.service.List.Get(user_id, list_id)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"list": list,
	})
}

func (h *Handler) getAllList(c *gin.Context) {
	user_id := c.GetInt("user_id")
	if user_id == 0 {
		h.newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	lists, err := h.service.List.GetAll(user_id)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"lists": lists,
	})
}

func (h *Handler) updateList(c *gin.Context) {
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
	var input domain.UpdateList

	if err := c.Bind(&input); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.List.Update(user_id, list_id, input.Title, input.Description); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"status": "updated",
	})

}

func (h *Handler) deleteList(c *gin.Context) {
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
	if err := h.service.List.Delete(user_id, list_id); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"status": "deleted",
	})

}
