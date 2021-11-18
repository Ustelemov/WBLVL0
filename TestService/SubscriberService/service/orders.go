package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/ustelemov/WBLVL0/TestService/cache"
	"github.com/ustelemov/WBLVL0/TestService/event"
	"github.com/ustelemov/WBLVL0/TestService/repository"
	"github.com/ustelemov/WBLVL0/TestService/schema"
)

type OrdersService struct {
	repo  repository.Repository
	Cache cache.RepositoryCache
}

func NewOrdersService(repo repository.Repository, cache cache.RepositoryCache) *OrdersService {
	return &OrdersService{repo: repo, Cache: cache}
}

func (service *OrdersService) GetOrderByUUID(uuid string) (*schema.Order, error) {
	return service.Cache.GetOrderByUUID(uuid)
}

func (service *OrdersService) SaveOrderInCache(ord *schema.Order) error {
	return service.Cache.SaveOrderInCache(ord)
}

func (service *OrdersService) SaveOrderInRepository(ord *schema.Order) error {
	return service.repo.SaveOrderInRepository(ord)
}

func (service *OrdersService) GetAllOrdersInCache() error {
	ordersJsonArr, err := service.repo.GetAllOrders()

	if err != nil {
		return fmt.Errorf("failed on load all orders in cache: %v", err.Error())
	}

	maps := make(map[string]*schema.Order)

	for _, v := range ordersJsonArr {
		uuid := v.OrderUID
		order := new(schema.Order)
		err := json.Unmarshal(v.JSON, order)

		if err != nil {
			return fmt.Errorf("cannot unmarshal JSON to Order")
		}
		maps[uuid] = order

	}

	service.Cache.ChangeMapRepositoryCache(maps)

	return nil
}

func (service *OrdersService) ProccessOrderMessage(msg *event.OrderMessage) {

	isValid := validateMessage(msg)

	if !isValid {
		logrus.Errorf("message is not valid as Order")
		return
	}

	if err := service.SaveOrderInRepository(msg.Order); err != nil {
		logrus.Errorf("error while proccess order message: %s", err.Error())
		return
	}

	if err := service.SaveOrderInCache(msg.Order); err != nil {
		logrus.Errorf("error while proccess order message: %s", err.Error())
		return
	}
	logrus.Println("message with order-uuid: %s from: %s processed succesful", msg.Order.OrderUID, msg.CreatedAt)
}

func validateMessage(msg *event.OrderMessage) bool {
	if (msg.CreatedAt == time.Time{}) || (msg.Order == nil) || (msg.Order.OrderUID == "") {
		return false
	}
	return true
}

func UpdateCacheOnTringer(data []byte) {
	orderJSON := schema.OrderJSON{}

	err := json.Unmarshal(data, &orderJSON)

	if err != nil || orderJSON.OrderUID == "" {
		logrus.Printf("error while unmarshal JSON from database: %s", err.Error())
	}

}
