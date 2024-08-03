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
	FlagSendAddr   string `json:"address"`
	shaKey         string
	ReportInterval int `json:"report_interval"`
	PollInterval   int `json:"poll_interval"`
	maxRetries     int
	CryptoKey      string `json:"crypto_key"`
}

func NewConfig() (*Config, error) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	c := &Config{
		FlagSendAddr:   "localhost:8080",
		ReportInterval: 10,
		PollInterval:   2,
		maxRetries:     5,
		shaKey:         "",
		CryptoKey:      "",
	}

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
		c.FlagSendAddr = v
	}

	var err error
	if v, ok := os.LookupEnv("REPORT_INTERVAL"); v != "" && ok {
		if c.ReportInterval, err = strconv.Atoi(v); err != nil {
			return fmt.Errorf("ENV REPORT_INTERVAL: %s", err)
		}
	}
	if v, ok := os.LookupEnv("POLL_INTERVAL"); v != "" && ok {
		if c.PollInterval, err = strconv.Atoi(v); err != nil {
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
		c.CryptoKey = v
	}

	flag.StringVar(&configFile, "c", configFile, "Path to the JSON config file")
	flag.StringVar(&configFile, "config", configFile, "Path to the JSON config file")
	flag.StringVar(&c.CryptoKey, "crypto-key", c.CryptoKey, "path to the public key file")
	flag.StringVar(&c.FlagSendAddr, "a", c.FlagSendAddr, "address and port to run server")
	flag.IntVar(&c.ReportInterval, "r", c.ReportInterval, "reportInterval description")
	flag.IntVar(&c.PollInterval, "p", c.PollInterval, "pollInterval description")
	flag.IntVar(&c.maxRetries, "i", c.maxRetries, "maxRetries description")
	flag.StringVar(&c.shaKey, "k", c.shaKey, "key description")
	flag.Parse()

	if flag.NArg() > 0 {
		flag.PrintDefaults()
		return fmt.Errorf("%s", flag.Arg(0))
	}

	if configFile != "" {
		file, err := os.Open(configFile)
		if err != nil {
			log.Fatalf("Error opening config file: %v", err)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&c); err != nil {
			log.Fatalf("Error decoding config file: %v", err)
		}
	}

	flag.Parse()

	log.Println("agent: " + c.shaKey)
	log.Println("agent: " + c.FlagSendAddr)
	log.Println("agent: " + strconv.Itoa(c.ReportInterval))
	log.Println("agent: " + strconv.Itoa(c.PollInterval))
	log.Println("agent: " + strconv.Itoa(c.maxRetries))

	return nil
}
