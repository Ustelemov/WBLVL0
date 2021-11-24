package schema

type Order struct {
	OrderUID          string  `json:"order_uid" fake:"{uuid}"`
	Entry             string  `json:"entry" fake:"{lettern:10}"`
	InternalSignature string  `json:"internal_signature" fake:"{digitn:12}"`
	Payment           Payment `json:"payment"`
	Items             []Item  `json:"items" fakesize:"3"`
	Locale            string  `json:"locale" fake:"{randomstring:[ru UTF-8, en UTF-8]}"`
	CustomerID        string  `json:"customer_id" fake:"{uuid}"`
	TrackNumber       string  `json:"track_number" fake:"{digitn:16}"`
	DeliveryService   string  `json:"delivery_service" fake:"{randomstring:[Postal,DPD,BoxBerry]}"`
	Shardkey          string  `json:"shardkey" fake:"{lettern:12}"`
	SmID              int     `json:"sm_id" fake:"{number:0}"`
}

type Payment struct {
	Transaction  string `json:"transaction" fake:"{uuid}"`
	Currency     string `json:"currency" fake:"{currencyshort}"`
	Provider     string `json:"provider" fake:"{randomstring: CC}"`
	Amount       int    `json:"amount" fake:"skip"`
	PaymentDt    int    `json:"payment_dt" fake:"{number:1636088400,1636952400}"` //from 10:00 5.11.21 to 10:00 15.11.21
	Bank         string `json:"bank" fake:"{randomstring:[Sber,Tinkoff,Alfa]}"`
	DeliveryCost int    `json:"delivery_cost" fake:"{number:0, 2000}"`
	GoodsTotal   int    `json:"goods_total" fake:"{number:3,3}"` //same as len([]items) in Order Struct`
}

type Item struct {
	ChrtID     int    `json:"chrt_id" fake:"{number: 0, 1000}"`
	Price      int    `json:"price" fake:"{number: 1000, 20000}"`
	Rid        string `json:"rid" fake:"{uuid}"`
	Name       string `json:"name" fake:"{fruit}"`
	Sale       int    `json:"sale" fake:"{number: 0, 99}"`
	Size       string `json:"size" fake:"{number: 0, 10}"`
	TotalPrice int    `json:"total_price" fake:"skip"`
	NmID       int    `json:"nm_id" fake:"{digitn:14}"`
	Brand      string `json:"brand" fake:"{randomstring:[Russia, Belarus, Turkey, Tailand]}"`
}

type OrderOut struct {
	OrderUID        string `json:"order_uid"`
	Entry           string `json:"entry"`
	TotalPrice      int    `json:"total_price"`
	CustomerID      string `json:"customer_id"`
	TrackNumber     string `json:"track_number"`
	DeliveryService string `json:"delivery_service"`
}

type OrderJSON struct {
	OrderUID string `json:"uuid" db:"uuid"`
	JSON     []byte `json:"order_data" db:"order_data"`
}
