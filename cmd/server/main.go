package main

import (
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/services/server"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"log"
	"os"
)

func handleError(err error, message string) {
	if err != nil {
		log.Println(message + err.Error())
		os.Exit(1)
	}
}

func main() {
	s := storages.NewMemStorage()
	s.AddMetric("gauge", metrics.NewGauge())
	s.AddMetric("counter", metrics.NewCounter())

	if err := parseFlags(); err != nil {
		handleError(err, "error occurred while parsing flags: ")
	}
	if err := server.New(s).Run(flagRunAddr); err != nil {
		handleError(err, "error occurred while running http server: ")
	}

	log.Println("server started on " + flagRunAddr)
	os.Exit(0)
}
