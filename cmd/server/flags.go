package main

import (
	"flag"
	"fmt"
	"os"
)

var flagRunAddr string

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Println("Unknown flag:", flag.Arg(0))
		flag.PrintDefaults()
		os.Exit(1)
	}

	if val, ok := os.LookupEnv("ADDRESS"); ok {
		flagRunAddr = val
	}
}
