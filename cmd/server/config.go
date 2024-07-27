package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	ServerAddress   string `json:"address"`
	StoreInterval   int    `json:"store_interval"`
	FileStoragePath string `json:"store_file"`
	Restore         bool   `json:"restore"`
	DataBaseDSN     string `json:"database_dsn"`
	shaKey          string
	migrationsDir   string
	CryptoKey       string `json:"crypto_key"`
}

func NewConfig() (*Config, error) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	c := &Config{
		ServerAddress:   "localhost:8080",
		StoreInterval:   300,
		FileStoragePath: "/tmp/metrics-db.json",
		Restore:         true,
		DataBaseDSN:     "postgres://postgres:root@localhost:5432",
		shaKey:          "",
		CryptoKey:       "",
	}

	projectDir, _ := os.Getwd()
	c.migrationsDir = projectDir + "/internal/storages/migrations"

	if err := c.parseFlags(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) parseFlags() error {
	configFile := ""
	if v, ok := os.LookupEnv("CONFIG"); v != "" && ok {
		configFile = v
	}

	if v, ok := os.LookupEnv("ADDRESS"); v != "" && ok {
		c.ServerAddress = v
	}

	var err error
	if v, ok := os.LookupEnv("STORE_INTERVAL"); v != "" && ok {
		if c.StoreInterval, err = strconv.Atoi(v); err != nil {
			return fmt.Errorf("ENV STORE_INTERVAL: %s", err)
		}
	}

	if v, ok := os.LookupEnv("FILE_STORAGE_PATH"); v != "" && ok {
		c.FileStoragePath = v
	}

	if v, ok := os.LookupEnv("RESTORE"); v != "" && ok {
		if c.Restore, err = strconv.ParseBool(v); err != nil {
			return fmt.Errorf("ENV RESTORE: %s", err)
		}
	}

	if v, ok := os.LookupEnv("DATABASE_DSN"); v != "" && ok {
		c.DataBaseDSN = v
	}

	if v, ok := os.LookupEnv("KEY"); v != "" && ok {
		c.shaKey = v
	}

	if v, ok := os.LookupEnv("CRYPTO_KEY"); v != "" && ok {
		c.CryptoKey = v
	}

	flag.StringVar(&configFile, "c", configFile, "Path to the JSON config file")
	flag.StringVar(&configFile, "config", configFile, "Path to the JSON config file")
	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "address and port to run server")
	flag.IntVar(&c.StoreInterval, "i", c.StoreInterval, "storeInterval")
	flag.StringVar(&c.FileStoragePath, "f", c.FileStoragePath, "fileStoragePath")
	flag.BoolVar(&c.Restore, "r", c.Restore, "restore")
	flag.StringVar(&c.DataBaseDSN, "d", c.DataBaseDSN, "databaseDSN")
	flag.StringVar(&c.shaKey, "k", c.shaKey, "shaKey")
	flag.StringVar(&c.CryptoKey, "crypto-key", c.CryptoKey, "path to the private key file")
	flag.Parse()

	if flag.NArg() > 0 {
		flag.PrintDefaults()
		return fmt.Errorf("%s", flag.Arg(0))
	}

	if configFile != "" {
		file, err := os.Open(configFile)
		if err != nil {
			return err
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&c); err != nil {
			return err
		}
	}

	flag.Parse()

	log.Println("server: " + c.shaKey)
	log.Println("server: " + c.DataBaseDSN)
	log.Println("server: " + c.FileStoragePath)
	log.Println("server: " + c.ServerAddress)
	log.Println("server: " + strconv.Itoa(c.StoreInterval))
	log.Println("server: " + strconv.FormatBool(c.Restore))

	return nil
}

func (c *Config) GetServerAddress() string {
	return c.ServerAddress
}

func (c *Config) GetStoreInterval() int {
	return c.StoreInterval
}

func (c *Config) GetFileStoragePath() string {
	return c.FileStoragePath
}

func (c *Config) GetRestore() bool {
	return c.Restore
}

func (c *Config) GetDataBaseDSN() string {
	return c.DataBaseDSN
}

func (c *Config) GetShaKey() string {
	return c.shaKey
}

func (c *Config) GetMigrationsDir() string {
	return c.migrationsDir
}

func (c *Config) GetCryptoKey() string {
	return c.CryptoKey
}
