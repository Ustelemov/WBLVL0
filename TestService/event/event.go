package event

type EventStorage interface {
	Close() error
	PublishOrder(string, *OrderMessage) error
	SubscribeOnOrders(func(*OrderMessage)) error
	QueueSubscribeOnOrders(func(*OrderMessage)) error
}
