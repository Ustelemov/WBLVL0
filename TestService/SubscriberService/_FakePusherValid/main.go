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
		logrus.Fatalf("error while init config: %s", err.Error())
	}

	nes, err := event.NewNatsJsonEventStore(event.Config{
		ConnecUrl:          viper.GetString("nats.url"),
		Max_reconnects:     viper.GetInt("nats.opts.max_reconnects"),
		Reconnect_wait_sec: viper.GetInt("nats.opts.reconnect_wait_sec"),
		Subject:            viper.GetString("nats.subject"),
		Subsc_queue:        viper.GetString("nats.subsc_queue"),
	})
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

		err = nes.PublishOrder(message)

		if err != nil {
			logrus.Fatal(err)
		}

		fmt.Printf("added %d valid-message with order-uuid: %s from: %s\n", i, order.OrderUID, message.CreatedAt.Format(time.RFC1123))
		i += 1

		time.Sleep(20 * time.Second)

	}

}
