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

func parseFlags() {
	flag.StringVar(&flagSendAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "reportInterval description")
	flag.IntVar(&pollInterval, "p", 2, "pollInterval description")
	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Println("Unknown flag:", flag.Arg(0))
		flag.PrintDefaults()
		os.Exit(1)
	}

	if val, ok := os.LookupEnv("ADDRESS"); ok {
		flagSendAddr = val
	}
	if val, ok := os.LookupEnv("REPORT_INTERVAL"); ok {
		reportInterval, _ = strconv.Atoi(val)
	}
	if val, ok := os.LookupEnv("POLL_INTERVAL"); ok {
		pollInterval, _ = strconv.Atoi(val)
	}
}
