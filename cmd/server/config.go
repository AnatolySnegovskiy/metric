package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	flagRunAddr string
}

func NewConfig() (*Config, error) {
	c := &Config{
		flagRunAddr: "localhost:8080",
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

	flag.StringVar(&c.flagRunAddr, "a", c.flagRunAddr, "address and port to run server")
	flag.Parse()

	if flag.NArg() > 0 {

		flag.PrintDefaults()
		return fmt.Errorf("%s", flag.Arg(0))
	}

	return nil
}
