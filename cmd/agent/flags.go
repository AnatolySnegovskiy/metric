package main

import (
	"flag"
	"fmt"
	"os"
)

var flagSendAddr string
var reportInterval int
var pollInterval int

func parseFlags() {
	flag.StringVar(&flagSendAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "reportInterval description")
	flag.IntVar(&pollInterval, "p", 2, "pollInterval description")
	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Println("Unknown flag:", flag.Arg(0))
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println("flagSendAddr:", flagSendAddr)
	fmt.Println("pollInterval:", pollInterval)
	fmt.Println("reportInterval:", reportInterval)
}
