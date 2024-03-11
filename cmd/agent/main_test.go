package main

import (
	"bou.ke/monkey"
	"bytes"
	"errors"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func Test_Main(t *testing.T) {
	resetVars()
	s := storages.NewMemStorage()
	s.AddMetric("gauge", metrics.NewGauge())
	s.AddMetric("counter", metrics.NewCounter())
	r, w, _ := os.Pipe()
	os.Stdout = w

	quit := make(chan struct{})

	go func() {
		defer close(quit)
		go main()
		time.Sleep(1 * time.Second)
		_ = w.Close()
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		expectedOutput := "Agent started\n"
		assert.Contains(t, buf.String(), expectedOutput, "Unexpected output. Expected: %s, got: %s", expectedOutput, buf.String())
	}()
}

func TestHandleShutdownSignal(t *testing.T) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go handleShutdownSignal(quit)
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
	fakeExit := func(int) {
		panic("os.Exit called")
	}
	patch := monkey.Patch(os.Exit, fakeExit)
	defer patch.Unpatch()
	err := errors.New("mock error")
	assert.PanicsWithValue(t, "os.Exit called", func() { handleError(err) }, "os.Exit was not called")
}

func TestHandleErrorWithNil(t *testing.T) {
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	handleError(nil)
	assert.Empty(t, logOutput.String())
}
