package handler

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_service "github.com/ustelemov/WBLVL0/TestService/service/mocks"
)

func TestHandler_NewHandler(t *testing.T) {
	t.Run("CreatedOK", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		order := mock_service.NewMockOrders(c)

		handler := NewHandler(order)

		assert.NotNil(t, handler)
	})
}

func TestHandler_InitRoutes(t *testing.T) {
	t.Run("InitOK", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		order := mock_service.NewMockOrders(c)

		handler := NewHandler(order)
		eng := handler.InitRoutes()

		assert.NotNil(t, eng)
	})
}
