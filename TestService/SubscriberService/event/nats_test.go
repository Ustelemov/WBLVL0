package event

import (
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/ustelemov/WBLVL0/TestService/schema"
)

var DefaultTestOptions = server.Options{
	Host:                  "127.0.0.1",
	Port:                  4222,
	NoLog:                 true,
	NoSigs:                true,
	MaxControlLine:        4096,
	DisableShortFirstPing: true,
}

var testConfig = Config{
	ConnecUrl:          "nats://127.0.0.1:4222",
	Max_reconnects:     10,
	Reconnect_wait_sec: 100,
	Subject:            "test_subj",
	Subsc_queue:        "test_queue",
}

func RunDefaultServer() (*server.Server, error) {
	s, err := server.NewServer(&DefaultTestOptions)
	if err != nil || s == nil {
		return nil, fmt.Errorf("error while creating nats-server: %s", err.Error())
	}
	s.Start()
	return s, nil
}

func TestNatsEventStorage_setupConnOptions(t *testing.T) {
	s, err := RunDefaultServer()

	if err != nil {
		t.Errorf("error while running nats-server: %s", err.Error())
	}

	defer s.Shutdown()

	opts := setupConnOptions(testConfig)

	nes, err := nats.Connect(testConfig.ConnecUrl, opts...)

	if nes != nil {
		defer nes.Close()
	}

	t.Run("ConnectionOptionsOK", func(t *testing.T) {
		assert.NoError(t, err)
	})

}

func TestNatsEventStorage_NewNatsJsonEventStore(t *testing.T) {
	s, err := RunDefaultServer()

	if err != nil {
		t.Errorf("error while running nats-server: %s", err.Error())
	}

	defer s.Shutdown()

	nes, err := NewNatsJsonEventStore(testConfig)

	t.Run("CreatedOK", func(t *testing.T) {
		assert.NoError(t, err)
		assert.NotNil(t, nes)
	})

}

func TestNatsEventStorage_Close(t *testing.T) {
	s, err := RunDefaultServer()

	if err != nil {
		t.Errorf("error while running nats-server: %s", err.Error())
	}

	defer s.Shutdown()

	nes, err := NewNatsJsonEventStore(testConfig)

	if err != nil {
		t.Errorf(err.Error())
	}

	err = nes.Close()

	t.Run("ClosedOK", func(t *testing.T) {
		assert.NoError(t, err)
	})

}
func TestNatsEventStorage_PublishOrder(t *testing.T) {
	s, err := RunDefaultServer()

	if err != nil {
		t.Errorf("error while running nats-server: %s", err.Error())
		return
	}

	defer s.Shutdown()

	nes, err := NewNatsJsonEventStore(testConfig)

	if err != nil {
		t.Errorf(err.Error())
	}

	msg := &OrderMessage{CreatedAt: time.Now(), Order: new(schema.Order)}

	err = nes.PublishOrder(msg)

	t.Run("PublishOK", func(t *testing.T) {
		assert.NoError(t, err)
	})
}

func TestNatsEventStorage_QueueSubscribeOnOrders(t *testing.T) {
	s, err := RunDefaultServer()

	if err != nil {
		t.Errorf("error while running nats-server: %s", err.Error())
		return
	}

	defer s.Shutdown()

	nes, err := NewNatsJsonEventStore(testConfig)

	if err != nil {
		t.Errorf(err.Error())
	}

	err = nes.QueueSubscribeOnOrders(func(om *OrderMessage) {
		_ = om
	})

	t.Run("QueueSubscribeOK", func(t *testing.T) {
		assert.NoError(t, err)
	})
}

func TestNatsEventStorage_SubscribeOnOrders(t *testing.T) {
	s, err := RunDefaultServer()

	if err != nil {
		t.Errorf("error while running nats-server: %s", err.Error())
		return
	}

	defer s.Shutdown()

	nes, err := NewNatsJsonEventStore(testConfig)

	if err != nil {
		t.Errorf(err.Error())
	}

	err = nes.SubscribeOnOrders(func(om *OrderMessage) {
		_ = om
	})

	t.Run("SubscribeOK", func(t *testing.T) {
		assert.NoError(t, err)
	})
}
