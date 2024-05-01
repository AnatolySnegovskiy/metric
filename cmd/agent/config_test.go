package main

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {
	t.Run("ENV", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("ADDRESS", "127.0.0.1:8080")
		_ = os.Setenv("REPORT_INTERVAL", "20")
		_ = os.Setenv("POLL_INTERVAL", "5")
		_ = os.Setenv("KEY", "test")
		_ = os.Setenv("RATE_LIMIT", "222")
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "127.0.0.1:8080", config.flagSendAddr, "expected default address")
		assert.Equal(t, 20, config.reportInterval, "expected default report interval")
		assert.Equal(t, 5, config.pollInterval, "expected default poll interval")
		assert.Equal(t, "test", config.shaKey, "expected default poll interval")
		assert.Equal(t, 222, config.maxRetries, "expected default poll interval")
	})

	t.Run("ENV_EMPTY", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("ADDRESS", "")
		_ = os.Setenv("REPORT_INTERVAL", "")
		_ = os.Setenv("POLL_INTERVAL", "")
		_ = os.Setenv("KEY", "")
		_ = os.Setenv("RATE_LIMIT", "")
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "localhost:8080", config.flagSendAddr, "expected default address")
		assert.Equal(t, 10, config.reportInterval, "expected default report interval")
		assert.Equal(t, 2, config.pollInterval, "expected default poll interval")
		assert.Equal(t, "", config.shaKey, "expected default poll interval")
		assert.Equal(t, 5, config.maxRetries, "expected default poll interval")
	})

	t.Run("ENV_ERROR_REPORT_INTERVAL", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("ADDRESS", "127.0.0.1:8080")
		_ = os.Setenv("REPORT_INTERVAL", "Error")
		_ = os.Setenv("POLL_INTERVAL", "5")
		_, err := NewConfig()
		assert.Error(t, err)
	})

	t.Run("ENV_ERROR_POLL_INTERVAL", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("ADDRESS", "127.0.0.1:8080")
		_ = os.Setenv("REPORT_INTERVAL", "10")
		_ = os.Setenv("POLL_INTERVAL", "Error")
		_, err := NewConfig()
		assert.Error(t, err)
	})

	t.Run("ENV_ERROR_RATE_LIMIT", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("ADDRESS", "127.0.0.1:8080")
		_ = os.Setenv("REPORT_INTERVAL", "10")
		_ = os.Setenv("RATE_LIMIT", "Error")
		_, err := NewConfig()
		assert.Error(t, err)
	})

	t.Run("CMD", func(t *testing.T) {
		resetVars()
		os.Args = []string{"cmd", "-a=127.0.10.1:8080", "-r=15", "-p=66", "-k=1234", "-i=10"}
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "127.0.10.1:8080", config.flagSendAddr, "expected default address")
		assert.Equal(t, 15, config.reportInterval, "expected default report interval")
		assert.Equal(t, 66, config.pollInterval, "expected default poll interval")
		assert.Equal(t, "1234", config.shaKey, "expected default poll interval")
		assert.Equal(t, 10, config.maxRetries, "expected default poll interval")
	})

	t.Run("CMD_OVERRIDE_ENV", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("ADDRESS", "127.0.0.1:8080")
		_ = os.Setenv("REPORT_INTERVAL", "20")
		_ = os.Setenv("POLL_INTERVAL", "5")
		_ = os.Setenv("KEY", "1234")
		_ = os.Setenv("RATE_LIMIT", "10")

		os.Args = []string{"cmd", "-a=127.21.10.1:8080", "-r=100", "-p=500", "-k=key", "-i=25"}
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "127.21.10.1:8080", config.flagSendAddr, "expected default address")
		assert.Equal(t, 100, config.reportInterval, "expected default report interval")
		assert.Equal(t, 500, config.pollInterval, "expected default poll interval")
		assert.Equal(t, "key", config.shaKey, "expected default poll interval")
		assert.Equal(t, 25, config.maxRetries, "expected default poll interval")
	})

	t.Run("CMD_ERROR_FLAG", func(t *testing.T) {
		resetVars()
		os.Args = []string{"cmd", "-B=127.0.10.1:8080", "-r=15", "-p=66"}
		_, err := NewConfig()
		assert.Error(t, err)
	})

	t.Run("CMD_ERROR_VALUE_R", func(t *testing.T) {
		resetVars()
		os.Args = []string{"cmd", "-a=127.0.10.1:8080", "-r=Error", "-p=1"}
		_, err := NewConfig()
		assert.Error(t, err)
	})

	t.Run("CMD_ERROR_VALUE_I", func(t *testing.T) {
		resetVars()
		os.Args = []string{"cmd", "-a=127.0.10.1:8080", "-i=Error", "-p=1"}
		_, err := NewConfig()
		assert.Error(t, err)
	})

	t.Run("DefaultValues", func(t *testing.T) {
		resetVars()
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "localhost:8080", config.flagSendAddr, "expected default address")
		assert.Equal(t, 10, config.reportInterval, "expected default report interval")
		assert.Equal(t, 2, config.pollInterval, "expected default poll interval")
		assert.Equal(t, "", config.shaKey, "expected default poll interval")
		assert.Equal(t, 5, config.maxRetries, "expected default poll interval")
	})
}

func resetVars() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cmd"}
	os.Clearenv()
}
