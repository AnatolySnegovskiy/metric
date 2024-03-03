package main

import (
	"flag"
)

var flagSendAddr string
var reportInterval int
var pollInterval int

func parseFlags() {
	flag.StringVar(&flagSendAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, " Sets the report interval metrics (default is 10 seconds)")
	flag.IntVar(&pollInterval, "p", 2, "Sets the pollInterval for runtime metrics (default is 2 seconds)")
	flag.Parse()
}
