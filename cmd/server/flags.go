package main

import (
	"flag"
	"fmt"
	"os"
)

var flagRunAddr string

func parseFlags() {
	flagRunAddr := flag.String("a", ":8080", "address and port to run server")
	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Println("Unknown flag:", flag.Arg(0))
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println("flagSendAddr:", *flagRunAddr)
}
