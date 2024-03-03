package main

import (
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/server"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"log"
)

func main() {
	storage := storages.NewMemStorage()
	storage.AddMetric("gauge", metrics.NewGauge())
	storage.AddMetric("counter", metrics.NewCounter())

	parseFlags()
	log.Println("server started on " + flagRunAddr)
	s := server.New(storage)

	if err := s.Run(flagRunAddr); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
