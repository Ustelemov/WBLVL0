package event

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

type NatsEventStorage struct {
	Conn         *nats.EncodedConn
	Subsctiption *nats.Subscription
}

func NewNatsJsonEventStore(connectUrl string, f func(*nats.Conn, *nats.Subscription, error)) (*NatsEventStorage, error) {

	conn, err := nats.Connect(connectUrl, nats.ErrorHandler(f))

	if err != nil {
		return nil, fmt.Errorf("error when trying connect to nats: %s", err)
	}

	econn, err := nats.NewEncodedConn(conn, "json")

	if err != nil {
		return nil, fmt.Errorf("error when trying encode connection: %s", err)
	}

	return &NatsEventStorage{Conn: econn}, nil

}

func (nes *NatsEventStorage) Close() {
	if nes.Conn != nil {
		nes.Conn.Close()
	}

	if nes.Subsctiption != nil {
		nes.Subsctiption.Unsubscribe()
	}
}

func (nes *NatsEventStorage) PublishOrder(subject string, ord *OrderMessage) error {
	err := nes.Conn.Publish(subject, ord)

	if err != nil {
		return fmt.Errorf("error while publishing: %s", err)
	}

	return nil
}

func (nes *NatsEventStorage) SubscribeOnOrders(f func(*OrderMessage)) (err error) {
	m := OrderMessage{}

	nes.Subsctiption, err = nes.Conn.Subscribe(m.Key(), func(ord *OrderMessage) {
		f(ord)
	})

	if err != nil {
		return fmt.Errorf("error while subscribing: %s", err)

	}

	return nil
}
