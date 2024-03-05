package main

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/services/agent"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"log"
	"os"
)

func handleError(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	s := storages.NewMemStorage()
	s.AddMetric("gauge", metrics.NewGauge())
	s.AddMetric("counter", metrics.NewCounter())

	fmt.Println("Agent started")
	err := agent.New(s).Run()
	handleError(err)
}
