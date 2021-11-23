package cache

import "github.com/ustelemov/WBLVL0/TestService/schema"

type Orders interface {
	GetOrderByUUID(string) *schema.Order
	SaveOrder(*schema.Order) error
	ChangeMapRepositoryCache(map[string]*schema.Order)
	GetAllOrders() []*schema.Order
}

type RepositoryCache interface {
	Orders
}
