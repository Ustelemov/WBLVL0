package handler

import (
	"github.com/gin-contrib/cors"
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
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://10.20.220.35:3000"}
	corsConfig.AllowCredentials = true

	corsConfig.AddAllowMethods("OPTIONS", "GET")

	router.Use(cors.New(corsConfig))

	api := router.Group("/orders")
	{
		api.GET("", h.getOrders)
		api.GET("/order", h.getOrderByUUID)
	}
	router.GET("/status", h.getStatus)

	return router
}
