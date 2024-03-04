package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var flagSendAddr string
var reportInterval int
var pollInterval int

func parseFlags() error {
	flag.StringVar(&flagSendAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "reportInterval description")
	flag.IntVar(&pollInterval, "p", 2, "pollInterval description")
	flag.Parse()

	if flag.NArg() > 0 {
		flag.PrintDefaults()
		return fmt.Errorf("%s", flag.Arg(0))
	}

	if val, ok := os.LookupEnv("ADDRESS"); ok {
		flagSendAddr = val
	}

	if val, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		if v, err := strconv.Atoi(val); err != nil {
			return fmt.Errorf("ENV REPORT_INTERVAL: %s", err)
		} else {
			reportInterval = v
		}
	}
	if val, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		if v, err := strconv.Atoi(val); err != nil {
			return fmt.Errorf("ENV POLL_INTERVAL: %s", err)
		} else {
			pollInterval = v
		}
	}

	return nil
}
