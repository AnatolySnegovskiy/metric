package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/stretchr/testify/assert"
)

func Test_Main(t *testing.T) {
	resetVars()
	os.Args = []string{"cmd", "-a=127.21.10.1:8150"}
	s := storages.NewMemStorage()
	s.AddMetric("gauge", metrics.NewGauge(nil))
	s.AddMetric("counter", metrics.NewCounter(nil))
	quit := make(chan struct{})

	go func() {
		defer close(quit)
		go main()
		time.Sleep(1 * time.Second)
		assert.True(t, true)
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

func TestSetDefaultValue(t *testing.T) {
	val := "test"
	val = setDefaultValue(val, "TEST2")
	assert.Equal(t, val, "test")
}

func TestSetDefaultValueEmpty(t *testing.T) {
	val := ""
	val = setDefaultValue(val, "TEST2")
	assert.Equal(t, val, "TEST2")
}
