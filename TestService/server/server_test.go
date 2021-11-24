package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testPort = "8632"

func TestServer_Run(t *testing.T) {

	t.Run("RunOK", func(t *testing.T) {
		server := new(Server)
		var err error

		go func() {
			if err = server.Run(testPort, http.NotFoundHandler()); err != nil && err != http.ErrServerClosed {
				return
			}
		}()

		timer := time.NewTimer(1 * time.Second)
		<-timer.C

		assert.NoError(t, err)
		assert.NotNil(t, server.httpServer)

		server.Shutdown(context.Background())

	})
}

func TestServer_Shutdown(t *testing.T) {

	t.Run("ShutdownOK", func(t *testing.T) {
		server := new(Server)

		go func() {
			if err := server.Run(testPort, http.NotFoundHandler()); err != nil && err != http.ErrServerClosed {
				return
			}
		}()

		timer := time.NewTimer(1 * time.Second)
		<-timer.C

		err := server.Shutdown(context.Background())

		assert.NoError(t, err)

	})
}
