package main

import (
	"bou.ke/monkey"
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func Test_Main(t *testing.T) {
	os.Args = []string{"cmd", "-a=localhost:8150"}
	go func() {
		main()
	}()
	var stdoutBuf bytes.Buffer
	log.SetOutput(&stdoutBuf)
	time.Sleep(1 * time.Second)
	assert.Contains(t, stdoutBuf.String(), "server started on localhost:8150", "Expected start message not found in the console output")
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
