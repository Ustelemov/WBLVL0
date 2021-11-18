package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ustelemov/WBLVL0/TestService/service"
)

type Handler struct {
	service service.Service
}

func NewHandler(s service.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		api.GET("/:uuid", h.getOrderByUUID)
	}

	return router

}
