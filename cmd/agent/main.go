package main

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/services"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"time"
)

const pollInterval = 2
const reportInterval = 10

func main() {
	storage := storages.NewMemStorage()
	storage.AddMetric("gauge", metrics.NewGauge())
	storage.AddMetric("counter", metrics.NewCounter())

	updateTicker := time.Tick(pollInterval * time.Second)
	sendTicker := time.Tick(reportInterval * time.Second)

	go services.UpdateStoragePeriodically(updateTicker, storage)
	go services.SendMetricsPeriodically(sendTicker, storage)
	fmt.Println("Agent started")
	// Ждем завершения работы программы (например, через сигнал ОС)
	select {}
}
