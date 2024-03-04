package main

import (
	"flag"
	"fmt"
	"os"
)

var flagRunAddr string

func parseFlags() error {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if flag.NArg() > 0 {
		flag.PrintDefaults()
		return fmt.Errorf("%s", flag.Arg(0))
	}

	if val, ok := os.LookupEnv("ADDRESS"); ok {
		flagRunAddr = val
	}

	return nil
}
