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

func (service *OrdersService) GetOrderByUUID(uuid string) *schema.Order {
	return service.Cache.GetOrderByUUID(uuid)
}

func (service *OrdersService) SaveOrderInCache(ord *schema.Order) error {
	return service.Cache.SaveOrder(ord)
}

func (service *OrdersService) SaveOrderInRepository(ord *schema.Order) error {
	return service.repo.SaveOrderInRepository(ord)
}

func (service *OrdersService) LoadAllOrdersInCache() error {
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
			return fmt.Errorf("failed on load all orders in cache: cannot unmarshal JSON to Order")
		}
		maps[uuid] = order

	}

	service.Cache.ChangeMapRepositoryCache(maps)

	return nil
}

func (service *OrdersService) ProccessOrderMessage(msg *event.OrderMessage) {

	logrus.Println("Got new message on proccessing")

	isValid := validateMessage(msg)

	if !isValid {
		logrus.Errorf("error while proccess order message: message is not valid as Order")
		return
	}

	if err := service.SaveOrderInRepository(msg.Order); err != nil {
		logrus.Println("error while proccess order message: cannot save order in repository: %s", err.Error())
		return
	}

	if err := service.SaveOrderInCache(msg.Order); err != nil {
		logrus.Println("error while proccess order message: cannot save order in cache %s", err.Error())
		return
	}
	logrus.Println("succesful proccesed message: order-uuid- %s from- %s",
		msg.Order.OrderUID, msg.CreatedAt.Format(time.RFC1123))
}

func validateMessage(msg *event.OrderMessage) bool {
	if (msg.CreatedAt == time.Time{}) || (msg.Order == nil) || (msg.Order.OrderUID == "") {
		return false
	}
	return true
}

func (orderService *OrdersService) UpdateCacheOnTringer(data []byte) {

	order := schema.Order{}

	err := json.Unmarshal(data, &order)

	if err != nil {
		logrus.Printf("error while update cache on update-triger: cannot unmarshal JSON from database: %s", err.Error())
	}

	orderService.Cache.SaveOrder(&order)

}

func (orderService *OrdersService) GetOrderOutByUUID(uuid string) *schema.OrderOut {
	order := orderService.GetOrderByUUID(uuid)

	if order == nil {
		return nil
	}

	totalPrice := order.Payment.Amount + order.Payment.DeliveryCost

	orderOut := schema.OrderOut{
		OrderUID:        order.OrderUID,
		Entry:           order.Entry,
		TotalPrice:      totalPrice,
		CustomerID:      order.CustomerID,
		TrackNumber:     order.TrackNumber,
		DeliveryService: order.DeliveryService,
	}

	return &orderOut

}

func (orderService *OrdersService) GetAllUUIDsInCache() *[]string {
	orders := orderService.Cache.GetAllOrders()

	uuids := make([]string, len(orders))

	for k, v := range orders {
		uuids[k] = v.OrderUID
	}

	return &uuids
}
