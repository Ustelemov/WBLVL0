package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ustelemov/WBLVL0/TestService/schema"
)

func TestOrderMessage_CreateOrderMessage(t *testing.T) {
	ord := CreateOrderMessage(&schema.Order{OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b51"}, time.Now())

	t.Run("CreatedOK", func(t *testing.T) {
		assert.NotNil(t, ord)
	})
}

func TestOrderMessage_GetCreatedTime(t *testing.T) {
	time := time.Now()
	ord := CreateOrderMessage(&schema.Order{OrderUID: "91bf621d-aaf3-4a0e-9048-329720a49b51"}, time)

	t.Run("CreatedOK", func(t *testing.T) {
		message_time := ord.GetCreatedTime()
		assert.Equal(t, time, message_time)
	})
}
