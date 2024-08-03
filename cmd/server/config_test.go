package main

import (
	"flag"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("ENV", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("ADDRESS", "127.0.0.1:8080")
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "127.0.0.1:8080", config.GetServerAddress(), "expected default address")
	})

	t.Run("ENV_EMPTY", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("ADDRESS", "")
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "localhost:8080", config.GetServerAddress(), "expected default address")
	})

	t.Run("CMD", func(t *testing.T) {
		resetVars()
		os.Args = []string{"cmd", "-a=127.0.10.1:8080"}
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "127.0.10.1:8080", config.GetServerAddress(), "expected default address")
	})

	t.Run("CMD_OVERRIDE_ENV", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("ADDRESS", "127.0.0.1:8080")
		os.Args = []string{"cmd", "-a=127.21.10.1:8080"}
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "127.21.10.1:8080", config.GetServerAddress(), "expected default address")
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
		assert.Equal(t, "localhost:8080", config.GetServerAddress(), "expected default address")
	})

	t.Run("ENV_STORE_INTERVAL", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("STORE_INTERVAL", "600")
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, 600, config.GetStoreInterval(), "expected default store interval")
	})

	t.Run("ENV_STORE_INTERVAL_NO_INT", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("STORE_INTERVAL", "noint")
		_, err := NewConfig()
		assert.Error(t, err)
	})

	t.Run("ENV_FILE_STORAGE_PATH", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("FILE_STORAGE_PATH", "/tmp/metrics-db-test.json")
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "/tmp/metrics-db-test.json", config.GetFileStoragePath(), "expected custom file storage path")
	})

	t.Run("ENV_RESTORE", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("RESTORE", "false")
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, false, config.GetRestore(), "expected restore to be false")
	})

	t.Run("ENV_DATABASE_DSN", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("DATABASE_DSN", "postgres://postgres:root@localhost:1000/public")
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "postgres://postgres:root@localhost:1000/public", config.GetDataBaseDSN(), "expected restore to be false")
	})

	t.Run("ENV_RESTORE_NO_BOOL", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("RESTORE", "nobool")
		_, err := NewConfig()
		assert.Error(t, err)
	})

	t.Run("ENV_KEY", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("KEY", "test")
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "test", config.GetShaKey(), "expected default sha key")
	})

	t.Run("ENV_CRYPTO_KEY", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("CRYPTO_KEY", "test")
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "test", config.GetCryptoKey(), "expected default crypto key")
	})

	t.Run("ENV_TRUSTED_SUBNET", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("TRUSTED_SUBNET", "127.0.0.1/24")
		config, err := NewConfig()
		assert.NoError(t, err)
		_, trustedSubnet, _ := net.ParseCIDR("127.0.0.1/24")
		assert.Equal(t, *trustedSubnet, config.GetTrustedSubnet(), "expected restore to be false")
	})

	t.Run("ENV_TRUSTED_SUBNET_ERROR", func(t *testing.T) {
		resetVars()
		_ = os.Setenv("TRUSTED_SUBNET", "127.0")
		_, err := NewConfig()
		assert.Error(t, err)
	})

	t.Run("ENV_CONFIG_FILE", func(t *testing.T) {
		_ = os.WriteFile(
			"config.json",
			[]byte(`{
				"address": "localhost:8080",
				"restore": true,
				"store_interval": 1,
				"store_file": "/path/to/file.db",
				"database_dsn": "",
				"crypto_key": "/path/to/key.pem",
				"trusted_subnet": ""
			}`), 0644)
		resetVars()
		_ = os.Setenv("CONFIG", "config.json")
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "localhost:8080", config.GetServerAddress(), "expected default address")
		assert.Equal(t, true, config.GetRestore(), "expected restore to be true")
		assert.Equal(t, 1, config.GetStoreInterval(), "expected default store interval")
		assert.Equal(t, "/path/to/file.db", config.GetFileStoragePath(), "expected default file storage path")
		assert.Equal(t, "", config.GetDataBaseDSN(), "expected default database DSN")
		assert.Equal(t, "/path/to/key.pem", config.GetCryptoKey(), "expected default crypto key")
		os.Remove("config.json")
	})

	t.Run("CMD_STORE_INTERVAL", func(t *testing.T) {
		resetVars()
		os.Args = []string{"cmd", "-i=600"}
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, 600, config.GetStoreInterval(), "expected custom store interval")
	})

	t.Run("CMD_FILE_STORAGE_PATH", func(t *testing.T) {
		resetVars()
		os.Args = []string{"cmd", "-f=/tmp/metrics-db-test.json"}
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "/tmp/metrics-db-test.json", config.GetFileStoragePath(), "expected custom file storage path")
	})

	t.Run("CMD_RESTORE", func(t *testing.T) {
		resetVars()
		os.Args = []string{"cmd", "-r=false"}
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, false, config.GetRestore(), "expected restore to be false")
	})

	t.Run("CMD_DATABASE_DSN", func(t *testing.T) {
		resetVars()
		os.Args = []string{"cmd", "-d=postgres://postgres:test@localhost:1000/public"}
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "postgres://postgres:test@localhost:1000/public", config.GetDataBaseDSN(), "expected restore to be false")
	})

	t.Run("CMD_KEY", func(t *testing.T) {
		resetVars()
		os.Args = []string{"cmd", "-k=1234"}
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "1234", config.GetShaKey(), "expected restore to be false")
	})

	t.Run("CMD_CRYPTO_KEY", func(t *testing.T) {
		resetVars()
		os.Args = []string{"cmd", "-crypto-key=1234"}
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "1234", config.GetCryptoKey(), "expected restore to be false")
	})

	t.Run("CMD_TRUSTED_SUBNET", func(t *testing.T) {
		resetVars()
		os.Args = []string{"cmd", "-t=128.0.0.0/24"}
		config, err := NewConfig()
		assert.NoError(t, err)
		_, trustedSubnet, _ := net.ParseCIDR("128.0.0.0/24")
		assert.Equal(t, *trustedSubnet, config.GetTrustedSubnet(), "expected restore to be false")
	})

	t.Run("CMD_CONFIG_FILE", func(t *testing.T) {
		_ = os.WriteFile(
			"config.json",
			[]byte(`{
				"address": "localhost:8080",
				"restore": true,
				"store_interval": 1,
				"store_file": "/path/to/file.db",
				"database_dsn": "",
				"crypto_key": "/path/to/key.pem",
				"trusted_subnet": "128.0.0.0/24"
			}`), 0644)
		resetVars()
		os.Args = []string{"cmd", "-c=config.json"}
		config, err := NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "localhost:8080", config.GetServerAddress(), "expected restore to be false")
		assert.Equal(t, true, config.GetRestore(), "expected restore to be false")
		assert.Equal(t, 1, config.GetStoreInterval(), "expected restore to be false")
		assert.Equal(t, "/path/to/file.db", config.GetFileStoragePath(), "expected restore to be false")
		assert.Equal(t, "", config.GetDataBaseDSN(), "expected restore to be false")
		assert.Equal(t, "/path/to/key.pem", config.GetCryptoKey(), "expected restore to be false")
		_ = os.Remove("config.json")

		_ = os.WriteFile(
			"config2.json",
			[]byte(`{
				"address": "localhost:1234",
				"restore": false,
				"store_interval": 10,
				"store_file": "/path/to/file2.db",
				"database_dsn": "123111",
				"crypto_key": "/path/to/key2.pem"
			}`), 0644)

		resetVars()
		os.Args = []string{"cmd", "-a=localhost:1111", "-config=config2.json"}
		config, err = NewConfig()
		assert.NoError(t, err)
		assert.Equal(t, "localhost:1111", config.GetServerAddress(), "expected restore to be false")
		assert.Equal(t, false, config.GetRestore(), "expected restore to be false")
		assert.Equal(t, 10, config.GetStoreInterval(), "expected restore to be false")
		assert.Equal(t, "/path/to/file2.db", config.GetFileStoragePath(), "expected restore to be false")
		assert.Equal(t, "123111", config.GetDataBaseDSN(), "expected restore to be false")
		assert.Equal(t, "/path/to/key2.pem", config.GetCryptoKey(), "expected restore to be false")
		_ = os.Remove("config2.json")
	})

	t.Run("FILE_ERROR", func(t *testing.T) {
		_ = os.WriteFile(
			"config.json",
			[]byte(`
				"address": "localhost:8080",
				"restore": true,
				"store_interval": 1,
				"store_file": "/path/to/file.db",
				"database_dsn": ""
			`), 0644)
		resetVars()
		os.Args = []string{"cmd", "-c=config.json"}
		_, err := NewConfig()
		assert.Error(t, err)

		_ = os.Remove("config.json")

		resetVars()
		os.Args = []string{"cmd", "-c=config.json"}
		_, err = NewConfig()
		assert.Error(t, err)
	})
}

func resetVars() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{"cmd"}
	os.Clearenv()
}
