package event

import (
	"time"

	"github.com/spf13/viper"
	"github.com/ustelemov/WBLVL0/TestService/schema"
)

type Message interface {
	Key() string
}

type OrderMessage struct {
	CreatedAt time.Time
	Order     *schema.Order
}

func (ocm *OrderMessage) Key() string {
	return viper.GetString("nats_subject")
}

func CreateOrderMessage(ord *schema.Order, time time.Time) *OrderMessage {
	return &OrderMessage{CreatedAt: time, Order: ord}
}
