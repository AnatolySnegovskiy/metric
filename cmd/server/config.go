package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	serverAddress   string
	storeInterval   int
	fileStoragePath string
	restore         bool
	dataBaseDSN     string
	shaKey          string
	migrationsDir   string
	cryptoKey       string
}

func NewConfig() (*Config, error) {
	c := &Config{
		serverAddress:   "localhost:8080",
		storeInterval:   300,
		fileStoragePath: "/tmp/metrics-db.json",
		restore:         true,
		dataBaseDSN:     "postgres://postgres:root@localhost:5432",
		shaKey:          "",
		cryptoKey:       "",
	}

	projectDir, _ := os.Getwd()
	c.migrationsDir = projectDir + "/internal/storages/migrations"

	if err := c.parseFlags(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) parseFlags() error {
	if v, ok := os.LookupEnv("ADDRESS"); v != "" && ok {
		c.serverAddress = v
	}

	var err error
	if v, ok := os.LookupEnv("STORE_INTERVAL"); v != "" && ok {
		if c.storeInterval, err = strconv.Atoi(v); err != nil {
			return fmt.Errorf("ENV STORE_INTERVAL: %s", err)
		}
	}

	if v, ok := os.LookupEnv("FILE_STORAGE_PATH"); v != "" && ok {
		c.fileStoragePath = v
	}

	if v, ok := os.LookupEnv("RESTORE"); v != "" && ok {
		if c.restore, err = strconv.ParseBool(v); err != nil {
			return fmt.Errorf("ENV RESTORE: %s", err)
		}
	}

	if v, ok := os.LookupEnv("DATABASE_DSN"); v != "" && ok {
		c.dataBaseDSN = v
	}

	if v, ok := os.LookupEnv("KEY"); v != "" && ok {
		c.shaKey = v
	}

	if v, ok := os.LookupEnv("CRYPTO_KEY"); v != "" && ok {
		c.cryptoKey = v
	}

	flag.StringVar(&c.serverAddress, "a", c.serverAddress, "address and port to run server")
	flag.IntVar(&c.storeInterval, "i", c.storeInterval, "storeInterval")
	flag.StringVar(&c.fileStoragePath, "f", c.fileStoragePath, "fileStoragePath")
	flag.BoolVar(&c.restore, "r", c.restore, "restore")
	flag.StringVar(&c.dataBaseDSN, "d", c.dataBaseDSN, "databaseDSN")
	flag.StringVar(&c.shaKey, "k", c.shaKey, "shaKey")
	flag.StringVar(&c.cryptoKey, "crypto-key", c.cryptoKey, "path to the private key file")

	flag.Parse()

	if flag.NArg() > 0 {

		flag.PrintDefaults()
		return fmt.Errorf("%s", flag.Arg(0))
	}
	log.Println("server: " + c.shaKey)
	log.Println("server: " + c.dataBaseDSN)
	log.Println("server: " + c.fileStoragePath)
	log.Println("server: " + c.serverAddress)
	log.Println("server: " + strconv.Itoa(c.storeInterval))
	log.Println("server: " + strconv.FormatBool(c.restore))

	return nil
}

func (c *Config) GetServerAddress() string {
	return c.serverAddress
}

func (c *Config) GetStoreInterval() int {
	return c.storeInterval
}

func (c *Config) GetFileStoragePath() string {
	return c.fileStoragePath
}

func (c *Config) GetRestore() bool {
	return c.restore
}

func (c *Config) GetDataBaseDSN() string {
	return c.dataBaseDSN
}

func (c *Config) GetShaKey() string {
	return c.shaKey
}

func (c *Config) GetMigrationsDir() string {
	return c.migrationsDir
}

func (c *Config) GetCryptoKey() string {
	return c.cryptoKey
}
