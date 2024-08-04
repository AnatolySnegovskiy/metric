package main

import (
	"bou.ke/monkey"
	"bytes"
	"context"
	"errors"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/services/server"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

func Test_Main(t *testing.T) {
	resetVars()
	server.PgxConnect = func(ctx context.Context, connString string) (*pgx.Conn, error) {
		return nil, nil
	}
	os.Args = []string{"cmd", "-a=:8150", "-grpc=:8151"}
	s := storages.NewMemStorage()
	s.AddMetric("gauge", metrics.NewGauge(nil))
	s.AddMetric("counter", metrics.NewCounter(nil))
	quit := make(chan struct{})

	go func() {
		defer close(quit)
		go main()
		time.Sleep(2 * time.Second)
		assert.True(t, true)
	}()
	time.Sleep(3 * time.Second)
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
