package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ustelemov/WBLVL0/TestService/schema"
)

func TestMapRepositoryCache_NewMapRepositoryCache(t *testing.T) {
	t.Run("MapNotNil", func(t *testing.T) {
		m := NewMapRepositoryCache()
		assert.NotNil(t, m.maps)
	})
}

func TestMapRepositoryCache_ReplaceMap(t *testing.T) {
	t.Run("ReplaceOK", func(t *testing.T) {
		m := NewMapRepositoryCache()
		nmap := map[string]*schema.Order{
			"91bf621d-aaf3-4a0e-9048-329720a49b50": {OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b50"},
		}
		m.ReplaceMap(nmap)

		assert.Equal(t, nmap, m.maps)

	})
}

func TestMapRepositoryCache_GetOrderByUUID(t *testing.T) {

	testTable := []struct {
		name           string
		testMap        map[string]*schema.Order
		testingUUID    string
		expectedResult *schema.Order
	}{
		{
			name: "Exist",
			testMap: map[string]*schema.Order{
				"91bf621d-aaf3-4a0e-9048-329720a49b51": {OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b51"},
			},
			testingUUID:    "91bf621d-aaf3-4a0e-9048-329720a49b51",
			expectedResult: &schema.Order{OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b51"},
		},
		{
			name: "NoExist",
			testMap: map[string]*schema.Order{
				"91bf621d-aaf3-4a0e-9048-329720a49b51": {OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b51"},
			},
			testingUUID:    "aaf3621d-91bf-4a0e-9048-329720a49b51",
			expectedResult: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			m := MapRepositoryCache{}
			m.ReplaceMap(testCase.testMap)
			res := m.GetOrderByUUID(testCase.testingUUID)

			assert.Equal(t, res, testCase.expectedResult)
		})
	}

}

func TestMapRepositoryCache_SaveOrder(t *testing.T) {

	testTable := []struct {
		name         string
		testMap      map[string]*schema.Order
		testingOrder *schema.Order
		expectError  bool
	}{
		{
			name: "SaveCorrect",
			testMap: map[string]*schema.Order{
				"91bf621d-aaf3-4a0e-9048-329720a49b51": {OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b51"},
			},
			testingOrder: &schema.Order{OrderUID: "aaf3621d-91bf-4a0e-9048-329720a49b51"},
			expectError:  false,
		},
		{
			name: "AlreadyExists",
			testMap: map[string]*schema.Order{
				"91bf621d-aaf3-4a0e-9048-329720a49b51": {OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b51"},
			},
			testingOrder: &schema.Order{OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b51"},
			expectError:  true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			m := MapRepositoryCache{}
			m.ReplaceMap(testCase.testMap)
			err := m.SaveOrder(testCase.testingOrder)

			if testCase.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}

}

func TestMapRepositoryCache_GetAllOrders(t *testing.T) {

	testTable := []struct {
		name            string
		testMap         map[string]*schema.Order
		expectingResult []*schema.Order
	}{
		{
			name: "NonEmptyMap",
			testMap: map[string]*schema.Order{
				"91bf621d-aaf3-4a0e-9048-329720a49b51": {OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b51"},
				"aaf3621d-aaf3-4a0e-9048-329720a49b50": {OrderUID: "aaf3621d-aaf3-4a0e-9048-329720a49b50"},
			},
			expectingResult: []*schema.Order{
				{OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b51"},
				{OrderUID: "aaf3621d-aaf3-4a0e-9048-329720a49b50"},
			},
		},
		{
			name:            "EmptyMap",
			testMap:         make(map[string]*schema.Order, 0),
			expectingResult: make([]*schema.Order, 0),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			m := MapRepositoryCache{}
			m.ReplaceMap(testCase.testMap)
			res := m.GetAllOrders()

			assert.Equal(t, res, testCase.expectingResult)
		})
	}

}
