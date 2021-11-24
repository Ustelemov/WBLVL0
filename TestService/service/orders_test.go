package service

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ustelemov/WBLVL0/TestService/cache"
	"github.com/ustelemov/WBLVL0/TestService/event"
	"github.com/ustelemov/WBLVL0/TestService/repository"
	"github.com/ustelemov/WBLVL0/TestService/schema"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestOrderService_NewOrdersService(t *testing.T) {
	t.Run("CreatedOK", func(t *testing.T) {
		s := NewOrdersService(&repository.PostgresDB{}, &cache.MapRepositoryCache{})
		assert.NotNil(t, s.repo)
		assert.NotNil(t, s.Cache)
	})
}

func TestOrderService_LoadAllOrdersInCache(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		db, mock, err := sqlxmock.Newx()

		if err != nil {
			t.Fatalf("error while openning mock database connection: %s", err)
		}
		defer db.Close()

		c := cache.NewMapRepositoryCache()
		s := NewOrdersService(&repository.PostgresDB{DB: db}, c)

		rows := sqlxmock.NewRows([]string{"uuid", "order_data"}).AddRow("91bf621d-aaf3-4a0e-9048-329720a49b5091bf621d-aaf3-4a0e-9048-329720a49b50", []byte(`{"entry": "ZWgNTalluZ", "items": [{"rid": "80ee0c92-a2f8-4a49-befa-3fe024d179b5", "name": "Honeydew", "sale": 92, "size": "9", "brand": "Belarus", "nm_id": 37705739069486, "price": 6700, "chrt_id": 704, "total_price": 536}, {"rid": "ccfea2e5-f4f4-4369-a525-b95f3dd9cc18", "name": "Peach", "sale": 72, "size": "0", "brand": "Tailand", "nm_id": 41629006865217, "price": 3552, "chrt_id": 343, "total_price": 995}, {"rid": "48fc8716-8c79-49e1-abbf-68ec2a92b45e", "name": "Banana", "sale": 60, "size": "1", "brand": "Belarus", "nm_id": 57571191302774, "price": 10652, "chrt_id": 793, "total_price": 4261}], "sm_id": 2012494407, "locale": "en UTF-8", "payment": {"bank": "Tinkoff", "amount": 5792, "currency": "XOF", "provider": "CC", "payment_dt": 1636945459, "goods_total": 3, "transaction": "f7ef43ae-3978-45d1-9fa2-03f3f90b3d16", "delivery_cost": 1269}, "shardkey": "oWcJVvmYOeLI", "order_uid": "91bf621d-aaf3-4a0e-9048-329720a49b50", "customer_id": "9a4560c8-28f2-4d10-bea9-17e660293d16", "track_number": "6089878449009872", "delivery_service": "BoxBerry", "internal_signature": "969536567842"}`))
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM ordersjson")).WillReturnRows(rows)

		res := s.LoadAllOrdersInCache()

		assert.NoError(t, res)
		assert.NotNil(t, s.Cache.GetOrderByUUID("91bf621d-aaf3-4a0e-9048-329720a49b5091bf621d-aaf3-4a0e-9048-329720a49b50"))

	})
}

func TestOrderService_validateMessage(t *testing.T) {
	testTable := []struct {
		name           string
		inputMessage   *event.OrderMessage
		expectedResult bool
	}{
		{
			name:           "Valid",
			inputMessage:   &event.OrderMessage{CreatedAt: time.Now(), Order: &schema.Order{OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b51"}},
			expectedResult: true,
		},
		{
			name:           "inValid",
			inputMessage:   &event.OrderMessage{CreatedAt: time.Now()},
			expectedResult: false,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			res := validateMessage(testCase.inputMessage)
			assert.Equal(t, testCase.expectedResult, res)

		})
	}
}

func TestOrderSerbice_GetOrderByUUID(t *testing.T) {
	testTable := []struct {
		name          string
		inputMap      map[string]*schema.Order
		searchUUID    string
		expectedOrder *schema.Order
	}{
		{
			name: "Exists",
			inputMap: map[string]*schema.Order{
				"91bf621d-aaf3-4a0e-9048-329720a49b51": {OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b51"},
				"aaf3621d-aaf3-4a0e-9048-329720a49b50": {OrderUID: "aaf3621d-aaf3-4a0e-9048-329720a49b50"},
			},
			searchUUID:    "91bf621d-aaf3-4a0e-9048-329720a49b51",
			expectedOrder: &schema.Order{OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b51"},
		},
		{
			name:          "Empty",
			inputMap:      make(map[string]*schema.Order, 0),
			searchUUID:    "91bf621d-aaf3-4a0e-9048-329720a49b51",
			expectedOrder: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := cache.NewMapRepositoryCache()
			c.ReplaceMap(testCase.inputMap)
			s := NewOrdersService(nil, c)
			resOrder := s.GetOrderByUUID(testCase.searchUUID)

			assert.Equal(t, testCase.expectedOrder, resOrder)
		})
	}
}

func TestOrderSerbice_GetOrderOutByUUID(t *testing.T) {
	testTable := []struct {
		name             string
		inputMap         map[string]*schema.Order
		searchUUID       string
		expectedOrderOut *schema.OrderOut
	}{
		{
			name: "Exists",
			inputMap: map[string]*schema.Order{
				"91bf621d-aaf3-4a0e-9048-329720a49b51": {OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b51"},
				"aaf3621d-aaf3-4a0e-9048-329720a49b50": {OrderUID: "aaf3621d-aaf3-4a0e-9048-329720a49b50"},
			},
			searchUUID:       "91bf621d-aaf3-4a0e-9048-329720a49b51",
			expectedOrderOut: &schema.OrderOut{OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b51"},
		},
		{
			name:             "Empty",
			inputMap:         make(map[string]*schema.Order, 0),
			searchUUID:       "91bf621d-aaf3-4a0e-9048-329720a49b51",
			expectedOrderOut: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := cache.NewMapRepositoryCache()
			c.ReplaceMap(testCase.inputMap)
			s := NewOrdersService(nil, c)
			resOrderOut := s.GetOrderOutByUUID(testCase.searchUUID)

			assert.Equal(t, testCase.expectedOrderOut, resOrderOut)

		})
	}
}

func TestOrderService_GetAllUUIDsInCache(t *testing.T) {

	testTable := []struct {
		name          string
		inputMap      map[string]*schema.Order
		expectedUUIDS []string
	}{
		{
			name: "NotEmpty",
			inputMap: map[string]*schema.Order{
				"aaf3621d-aaf3-4a0e-9048-329720a49b50": {OrderUID: "aaf3621d-aaf3-4a0e-9048-329720a49b50"},
			},
			expectedUUIDS: []string{
				"aaf3621d-aaf3-4a0e-9048-329720a49b50",
			},
		},
		{
			name:          "Empty",
			inputMap:      make(map[string]*schema.Order, 0),
			expectedUUIDS: make([]string, 0),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := cache.NewMapRepositoryCache()
			c.ReplaceMap(testCase.inputMap)
			s := NewOrdersService(nil, c)
			resUUIDS := s.GetAllUUIDsInCache()

			assert.Equal(t, resUUIDS, testCase.expectedUUIDS)

		})
	}

}
