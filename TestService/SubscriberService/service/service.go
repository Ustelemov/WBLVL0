package service

import (
	"github.com/ustelemov/WBLVL0/TestService/event"
	"github.com/ustelemov/WBLVL0/TestService/schema"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Orders interface {
	GetOrderByUUID(uuid string) *schema.Order
	SaveOrderInCache(*schema.Order) error
	SaveOrderInRepository(*schema.Order) error
	LoadAllOrdersInCache() error
	ProccessOrderMessage(*event.OrderMessage)
	UpdateCacheOnTringer([]byte)
	GetOrderOutByUUID(string) *schema.OrderOut
	GetAllUUIDsInCache() *[]string
}

type Service interface {
	Orders
}
