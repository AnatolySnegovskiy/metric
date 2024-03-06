package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/services/agent"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"testing"
)

func TestAgent_Run(t *testing.T) {
	resetVars()
	s := storages.NewMemStorage()
	s.AddMetric("gauge", metrics.NewGauge())
	s.AddMetric("counter", metrics.NewCounter())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go handleShutdownSignal(quit)

	c, err := NewConfig()
	if err != nil {
		t.Fatalf("failed to create config: %v", err)
	}

	a := agent.New(
		agent.Options{
			Storage:        s,
			PollInterval:   c.pollInterval,
			ReportInterval: c.reportInterval,
			SendAddr:       c.flagSendAddr,
		},
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		err := a.Run(ctx)
		assert.NoError(t, err)
	}()
}

func TestHandleNoError(t *testing.T) {
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

func TestHandleError(t *testing.T) {
	resetVars()
	t.Run("Error case", func(t *testing.T) {
		if os.Getenv("BE_CRASHER") == "1" {
			err := errors.New("mock error")
			var logOutput bytes.Buffer
			log.SetOutput(&logOutput)
			defer func() {
				log.SetOutput(os.Stderr)
			}()
			handleError(err)
		}

		cmd := exec.Command(os.Args[0], "-test.run=TestHandleError")
		cmd.Env = append(os.Environ(), "BE_CRASHER=1")
		err := cmd.Run()
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			return
		}
		assert.Contains(t, err.Error(), "mock error")
	})
}
