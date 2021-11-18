package cache

import "github.com/ustelemov/WBLVL0/TestService/schema"

type Orders interface {
	GetOrderByUUID(string) (*schema.Order, error)
	SaveOrderInCache(*schema.Order) error
	ChangeMapRepositoryCache(map[string]*schema.Order)
	GetAllOrders() ([]*schema.Order, error)
}

type RepositoryCache interface {
	Orders
}
