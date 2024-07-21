package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	flagSendAddr   string
	shaKey         string
	reportInterval int
	pollInterval   int
	maxRetries     int
	cryptoKey      string
}

func NewConfig() (*Config, error) {
	c := &Config{
		flagSendAddr:   "localhost:8080",
		reportInterval: 10,
		pollInterval:   2,
		maxRetries:     5,
		shaKey:         "",
		cryptoKey:      "",
	}

	if err := c.parseFlags(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) parseFlags() error {
	if v, ok := os.LookupEnv("ADDRESS"); v != "" && ok {
		c.flagSendAddr = v
	}

	var err error
	if v, ok := os.LookupEnv("REPORT_INTERVAL"); v != "" && ok {
		if c.reportInterval, err = strconv.Atoi(v); err != nil {
			return fmt.Errorf("ENV REPORT_INTERVAL: %s", err)
		}
	}
	if v, ok := os.LookupEnv("POLL_INTERVAL"); v != "" && ok {
		if c.pollInterval, err = strconv.Atoi(v); err != nil {
			return fmt.Errorf("ENV POLL_INTERVAL: %s", err)
		}
	}
	if v, ok := os.LookupEnv("RATE_LIMIT"); v != "" && ok {
		if c.maxRetries, err = strconv.Atoi(v); err != nil {
			return fmt.Errorf("ENV RATE_LIMIT: %s", err)
		}
	}
	if v, ok := os.LookupEnv("KEY"); v != "" && ok {
		c.shaKey = v
	}
	if v, ok := os.LookupEnv("CRYPTO_KEY"); v != "" && ok {
		c.cryptoKey = v
	}

	flag.StringVar(&c.cryptoKey, "crypto-key", c.cryptoKey, "path to the public key file")
	flag.StringVar(&c.flagSendAddr, "a", c.flagSendAddr, "address and port to run server")
	flag.IntVar(&c.reportInterval, "r", c.reportInterval, "reportInterval description")
	flag.IntVar(&c.pollInterval, "p", c.pollInterval, "pollInterval description")
	flag.IntVar(&c.maxRetries, "i", c.maxRetries, "maxRetries description")
	flag.StringVar(&c.shaKey, "k", c.shaKey, "key description")
	flag.Parse()

	if flag.NArg() > 0 {
		flag.PrintDefaults()
		return fmt.Errorf("%s", flag.Arg(0))
	}

	log.Println("agent: " + c.shaKey)
	log.Println("agent: " + c.flagSendAddr)
	log.Println("agent: " + strconv.Itoa(c.reportInterval))
	log.Println("agent: " + strconv.Itoa(c.pollInterval))
	log.Println("agent: " + strconv.Itoa(c.maxRetries))

	return nil
}
