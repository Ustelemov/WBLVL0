package event

import (
	"time"

	"github.com/ustelemov/WBLVL0/TestService/schema"
)

type Message interface {
	GetCreatedTime() time.Time
}

type OrderMessage struct {
	CreatedAt time.Time
	Order     *schema.Order
}

func (ocm *OrderMessage) GetCreatedTime() time.Time {
	return ocm.CreatedAt
}

func CreateOrderMessage(ord *schema.Order, time time.Time) *OrderMessage {
	return &OrderMessage{CreatedAt: time, Order: ord}
}
