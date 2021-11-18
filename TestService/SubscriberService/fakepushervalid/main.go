package main

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/ustelemov/WBLVL0/TestService/event"
	"github.com/ustelemov/WBLVL0/TestService/schema"
)

func CreateFakeOrder() (*schema.Order, error) {

	fake_order := schema.Order{}

	gofakeit.Struct(&fake_order)

	amount := 0 // sum total prices of items

	for i, v := range fake_order.Items {
		sale := float64(v.Sale*v.Price) / float64(100)
		totalPrice := v.Price - int(sale+0.5)

		fake_order.Items[i].TotalPrice = totalPrice

		amount += totalPrice
	}

	fake_order.Payment.Amount = amount

	return &fake_order, nil
}

func initConfig() error {
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {

	if err := initConfig(); err != nil {
		logrus.Fatalf("Error while init config: %s", err.Error())
	}

	nes, err := event.NewNatsJsonEventStore(viper.GetString("nats_url"), nil)

	if err != nil {
		logrus.Fatalf("error while creating nats-connection: %s", err.Error())
	}

	defer nes.Close()

	i := 1

	for {
		order, err := CreateFakeOrder()

		if err != nil {
			logrus.Fatal("error while creating fake order: %s", err)
		}

		message := event.CreateOrderMessage(order, time.Now())

		err = nes.PublishOrder(viper.GetString("nats_subject"), message)

		if err != nil {
			logrus.Fatal("error while publish valid fake message: %s", err)
		}

		fmt.Printf("Added %d valid-message with order-uuid: %s from: %s\n", i, order.OrderUID, message.CreatedAt)
		i += 1

		time.Sleep(20 * time.Second)

	}

}
