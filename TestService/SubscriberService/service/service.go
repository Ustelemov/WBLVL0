package service

import (
	"github.com/ustelemov/WBLVL0/TestService/event"
	"github.com/ustelemov/WBLVL0/TestService/schema"
)

type Orders interface {
	GetOrderByUUID(uuid string) (*schema.Order, error)
	SaveOrderInCache(*schema.Order) error
	SaveOrderInRepository(*schema.Order) error
	GetAllOrdersInCache() error
	ProccessOrderMessage(*event.OrderMessage)
	UpdateCacheOnTringer([]byte)
}

type Service interface {
	Orders
}
