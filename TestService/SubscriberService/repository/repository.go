package repository

import "github.com/ustelemov/WBLVL0/TestService/schema"

type Orders interface {
	SaveOrderInRepository(*schema.Order) error
	GetAllOrders() ([]schema.OrderJSON, error)
	CheckExists(*schema.Order) (bool, error)
	RunListenerDeamon(string, func([]byte))
}

type Repository interface {
	Orders
}
