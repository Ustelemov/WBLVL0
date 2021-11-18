package cache

import (
	"fmt"

	"github.com/ustelemov/WBLVL0/TestService/schema"
)

type MapRepositoryCache struct {
	maps map[string]*schema.Order
}

func NewMapRepositoryCache() *MapRepositoryCache {
	return &MapRepositoryCache{maps: make(map[string]*schema.Order)}
}

func (cache *MapRepositoryCache) ChangeMapRepositoryCache(maps map[string]*schema.Order) {
	cache.maps = maps
}

func (cache *MapRepositoryCache) GetOrderByUUID(uuid string) (*schema.Order, error) {
	v, ok := cache.maps[uuid]

	if !ok {
		return nil, fmt.Errorf("no uuid: %s found in cache", uuid)
	}

	return v, nil
}

func (cache *MapRepositoryCache) SaveOrderInCache(ord *schema.Order) error {
	uuid := ord.OrderUID

	if _, ok := cache.maps[uuid]; ok {
		return fmt.Errorf("uuid: %s already in cache", uuid)
	}

	cache.maps[uuid] = ord

	return nil
}

func (cache *MapRepositoryCache) GetAllOrders() ([]*schema.Order, error) {
	ordersArr := make([]*schema.Order, 0)

	for _, v := range cache.maps {
		ordersArr = append(ordersArr, v)
	}
	return ordersArr, nil
}
