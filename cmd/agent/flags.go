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
	flagSendAddr := flag.String("a", ":8080", "address and port to run server")
	pollInterval := flag.Int("p", 2, "pollInterval description")
	reportInterval := flag.Int("r", 3, "reportInterval description")

	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Println("Unknown flag:", flag.Arg(0))
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println("flagSendAddr:", *flagSendAddr)
	fmt.Println("pollInterval:", *pollInterval)
	fmt.Println("reportInterval:", *reportInterval)
}
