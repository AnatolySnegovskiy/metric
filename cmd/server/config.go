package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	flagRunAddr     string
	storeInterval   int
	fileStoragePath string
	restore         bool
	dataBaseDSN     string
}

func NewConfig() (*Config, error) {
	c := &Config{
		flagRunAddr:     "localhost:8080",
		storeInterval:   300,
		fileStoragePath: "/tmp/metrics-db.json",
		restore:         true,
		dataBaseDSN:     "postgres://postgres:root@localhost:5432/public",
	}

	if err := c.parseFlags(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) parseFlags() error {
	if val, ok := os.LookupEnv("ADDRESS"); val != "" && ok {
		c.flagRunAddr = val
	}

	var err error
	if val, ok := os.LookupEnv("STORE_INTERVAL"); val != "" && ok {
		if c.storeInterval, err = strconv.Atoi(val); err != nil {
			return fmt.Errorf("ENV STORE_INTERVAL: %s", err)
		}
	}

	if val, ok := os.LookupEnv("FILE_STORAGE_PATH"); val != "" && ok {
		c.fileStoragePath = val
	}

	if val, ok := os.LookupEnv("RESTORE"); val != "" && ok {
		if c.restore, err = strconv.ParseBool(val); err != nil {
			return fmt.Errorf("ENV RESTORE: %s", err)
		}
	}

	if val, ok := os.LookupEnv("DATABASE_DSN"); val != "" && ok {
		c.dataBaseDSN = val
	}

	flag.StringVar(&c.flagRunAddr, "a", c.flagRunAddr, "address and port to run server")
	flag.IntVar(&c.storeInterval, "i", c.storeInterval, "storeInterval")
	flag.StringVar(&c.fileStoragePath, "f", c.fileStoragePath, "fileStoragePath")
	flag.BoolVar(&c.restore, "r", c.restore, "restore")
	flag.StringVar(&c.dataBaseDSN, "d", c.dataBaseDSN, "databaseDSN")
	flag.Parse()

	if flag.NArg() > 0 {

		flag.PrintDefaults()
		return fmt.Errorf("%s", flag.Arg(0))
	}

	return nil
}
