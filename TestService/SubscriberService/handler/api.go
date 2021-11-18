package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getOrderByUUID(c *gin.Context) {
	var inputUUID string

	if err := c.BindJSON(inputUUID); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

}
