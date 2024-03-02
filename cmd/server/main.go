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

	s := server.New(storage)
	err := s.Run(`:8080`)

	if err != nil {
		log.Fatalf("start server: %v", err)
	}
}
