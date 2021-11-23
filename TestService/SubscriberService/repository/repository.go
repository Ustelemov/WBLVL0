package repository

import "github.com/ustelemov/WBLVL0/TestService/schema"

type Orders interface {
	SaveOrderInRepository(*schema.Order) error
	GetAllOrders() ([]schema.OrderJSON, error)
	CheckOrderExists(*schema.Order) (bool, error)
	RunListenerDeamon(string, func([]byte)) error
	Close()
}

type Repository interface {
	Orders
}
