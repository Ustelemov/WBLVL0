package event

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type NatsEventStorage struct {
	Conn         *nats.EncodedConn
	Subsctiption *nats.Subscription
	config       Config
}

type Config struct {
	ConnecUrl          string
	Max_reconnects     int
	Reconnect_wait_sec int
	Subject            string
	Subsc_queue        string
}

func NewNatsJsonEventStore(cfg Config) (*NatsEventStorage, error) {

	opts := setupConnOptions(cfg)

	conn, err := nats.Connect(cfg.ConnecUrl, opts...)

	if err != nil {
		return nil, fmt.Errorf("error when trying connect to nats: %s", err)
	}

	econn, err := nats.NewEncodedConn(conn, "json")

	if err != nil {
		return nil, fmt.Errorf("error when trying encode connection: %s", err)
	}

	return &NatsEventStorage{Conn: econn, config: cfg}, nil

}

func setupConnOptions(cfg Config) []nats.Option {

	opts := []nats.Option{}

	reconnectWait := time.Duration(cfg.Reconnect_wait_sec) * time.Second
	totalWait := reconnectWait * time.Duration(cfg.Max_reconnects)

	opts = append(opts, nats.ReconnectWait(reconnectWait))
	opts = append(opts, nats.MaxReconnects(cfg.Max_reconnects))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		logrus.Printf("nats disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		logrus.Printf("nats reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		logrus.Printf("nats exiting: %s", nc.LastError())
	}))
	opts = append(opts, nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
		logrus.Printf("nats subscribtion get error: %s ", err.Error())
	}))

	return opts
}

func (nes *NatsEventStorage) Close() error {
	if nes.Conn != nil {
		nes.Conn.Close()
	} else {
		return fmt.Errorf("error while trying to close nil nats-connection")
	}

	if nes.Subsctiption != nil {
		nes.Subsctiption.Unsubscribe()
	}

	return nil
}

func (nes *NatsEventStorage) PublishOrder(ord *OrderMessage) error {
	err := nes.Conn.Publish(nes.config.Subject, ord)

	if err != nil {
		return fmt.Errorf("error while publishing: %s", err)
	}

	return nil
}

func (nes *NatsEventStorage) QueueSubscribeOnOrders(f func(*OrderMessage)) (err error) {

	nes.Subsctiption, err = nes.Conn.QueueSubscribe(nes.config.Subject, nes.config.Subsc_queue, func(ord *OrderMessage) {
		f(ord)
	})

	if err != nil {
		return fmt.Errorf("error while queue-subscribing: %s", err)
	}

	nes.Conn.Flush()

	return nil
}

func (nes *NatsEventStorage) SubscribeOnOrders(f func(*OrderMessage)) (err error) {

	nes.Subsctiption, err = nes.Conn.Subscribe(nes.config.Subject, func(ord *OrderMessage) {
		f(ord)
	})

	if err != nil {
		return fmt.Errorf("error while subscribing: %s", err)

	}

	nes.Conn.Flush()

	return nil
}
