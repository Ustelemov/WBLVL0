package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/ustelemov/WBLVL0/TestService/schema"
	mock_service "github.com/ustelemov/WBLVL0/TestService/service/mocks"
)

func TestHandler_getOrderByUUID(t *testing.T) {
	type mockBehavior func(*mock_service.MockOrders, string)

	expectedOrder := &schema.OrderOut{
		OrderUID:        "91bf621d-aaf3-4a0e-9048-329720a49b50",
		Entry:           "ZWgNTalluZ",
		TotalPrice:      7061,
		CustomerID:      "9a4560c8-28f2-4d10-bea9-17e660293d16",
		TrackNumber:     "6089878449009872",
		DeliveryService: "BoxBerry",
	}

	testTable := []struct {
		name                string
		inputUUID           string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputUUID: "91bf621d-aaf3-4a0e-9048-329720a49b50",
			mockBehavior: func(s *mock_service.MockOrders, uuid string) {
				s.EXPECT().GetOrderOutByUUID(uuid).Return(expectedOrder)

			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"order_uid":"91bf621d-aaf3-4a0e-9048-329720a49b50","entry":"ZWgNTalluZ","total_price":7061,"customer_id":"9a4560c8-28f2-4d10-bea9-17e660293d16","track_number":"6089878449009872","delivery_service":"BoxBerry"}`,
		},
		{
			name:      "Fail",
			inputUUID: "sdfdsfsdf",
			mockBehavior: func(s *mock_service.MockOrders, uuid string) {
				s.EXPECT().GetOrderOutByUUID(uuid).Return(nil)

			},
			expectedStatusCode:  http.StatusNotFound,
			expectedRequestBody: `{"message":"Not found order"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			order := mock_service.NewMockOrders(c)
			testCase.mockBehavior(order, testCase.inputUUID)

			handler := NewHandler(order)

			r := gin.New()
			r.GET("/getbyuuid", handler.getOrderByUUID)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/getbyuuid", nil)
			q := req.URL.Query()
			q.Add("uuid", testCase.inputUUID)
			req.URL.RawQuery = q.Encode()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())

		})

	}

}

func TestHandler_getOrders(t *testing.T) {
	type mockBehavior func(*mock_service.MockOrders)

	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockOrders) {
				s.EXPECT().GetAllUUIDsInCache().Return([]string{"91bf621d-aaf3-4a0e-9048-329720a49b50"})
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"91bf621d-aaf3-4a0e-9048-329720a49b50"}`,
		},
		{
			name: "NotFound",
			mockBehavior: func(s *mock_service.MockOrders) {
				s.EXPECT().GetAllUUIDsInCache().Return(nil)
			},
			expectedStatusCode:  http.StatusNotFound,
			expectedRequestBody: `{"message":"Not found order"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			order := mock_service.NewMockOrders(c)
			testCase.mockBehavior(order)

			handler := NewHandler(order)

			r := gin.New()
			r.GET("/getorders", handler.getOrders)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/getorders", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})

	}

}
