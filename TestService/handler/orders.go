package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getOrderByUUID(c *gin.Context) {
	inputUUID := c.Query("uuid")

	orderOut := h.service.GetOrderOutByUUID(inputUUID)

	if orderOut == nil {
		newErrorResponse(c, http.StatusNotFound, "Not found order")
		return
	}

	c.JSON(http.StatusOK, orderOut)
}

func (h *Handler) getOrders(c *gin.Context) {
	orders := h.service.GetAllUUIDsInCache()

	if orders == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}

	c.JSON(http.StatusOK, orders)
}

func (h *Handler) getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, "OK")
}
