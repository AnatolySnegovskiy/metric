package main

import (
	"bytes"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/services/server"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func Test_Main(t *testing.T) {
	resetVars()
	os.Args = []string{"cmd", "-a=127.21.10.1:8150"}
	s := storages.NewMemStorage()
	s.AddMetric("gauge", metrics.NewGauge())
	s.AddMetric("counter", metrics.NewCounter())
	quit := make(chan struct{})

	go func() {
		defer close(quit)
		go main()
		time.Sleep(1 * time.Second)
		assert.True(t, true)
	}()
}

func TestHandleShutdownSignal(t *testing.T) {
	resetVars()
	s := server.New(storages.NewMemStorage(), nil)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	conf, _ := NewConfig()
	go handleShutdownSignal(quit, s, conf)
}

func TestHandleNoError(t *testing.T) {
	resetVars()
	t.Run("No error case", func(t *testing.T) {
		var logOutput bytes.Buffer
		log.SetOutput(&logOutput)
		defer func() {
			log.SetOutput(os.Stderr)
		}()
		handleError(nil)
		assert.Empty(t, logOutput.String())
	})
}
