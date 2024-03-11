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
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "127.0.0.1:8080", config.flagRunAddr, "expected default address")
	})

	t.Run("ENV_EMPTY", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("ADDRESS", "")
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "localhost:8080", config.flagRunAddr, "expected default address")
	})

	t.Run("CMD", func(t *testing.T) {
		resetVars()
		os.Args = []string{"cmd", "-a=127.0.10.1:8080"}
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "127.0.10.1:8080", config.flagRunAddr, "expected default address")
	})

	t.Run("CMD_OVERRIDE_ENV", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("ADDRESS", "127.0.0.1:8080")
		os.Args = []string{"cmd", "-a=127.21.10.1:8080"}
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "127.21.10.1:8080", config.flagRunAddr, "expected default address")
	})

	t.Run("CMD_ERROR_FLAG", func(t *testing.T) {
		resetVars()
		os.Args = []string{"cmd", "-B=127.0.10.1:8080", "-r=15"}
		_, err := NewConfig()
		assert.Error(t, err)
	})

	t.Run("DefaultValues", func(t *testing.T) {
		resetVars()
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "localhost:8080", config.flagRunAddr, "expected default address")
	})
}

func resetVars() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cmd"}
	os.Clearenv()
}
