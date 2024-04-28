package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	flagSendAddr   string
	shaKey         string
	reportInterval int
	pollInterval   int
	maxRetries     int
}

func NewConfig() (*Config, error) {
	c := &Config{
		flagSendAddr:   "localhost:8080",
		reportInterval: 10,
		pollInterval:   2,
		maxRetries:     5,
		shaKey:         "",
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
	if v, ok := os.LookupEnv("KEY"); v != "" && ok {
		c.shaKey = v
	}

	flag.StringVar(&c.flagSendAddr, "a", c.flagSendAddr, "address and port to run server")
	flag.IntVar(&c.reportInterval, "r", c.reportInterval, "reportInterval description")
	flag.IntVar(&c.pollInterval, "p", c.pollInterval, "pollInterval description")
	flag.StringVar(&c.shaKey, "k", c.shaKey, "pollInterval description")
	flag.Parse()

	if flag.NArg() > 0 {
		flag.PrintDefaults()
		return fmt.Errorf("%s", flag.Arg(0))
	}

	return nil
}
