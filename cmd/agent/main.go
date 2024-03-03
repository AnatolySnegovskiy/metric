package main

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/services"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"time"
)

func main() {
	storage := storages.NewMemStorage()
	storage.AddMetric("gauge", metrics.NewGauge())
	storage.AddMetric("counter", metrics.NewCounter())
	parseFlags()

	updateTicker := time.Tick(time.Duration(pollInterval) * time.Second)
	sendTicker := time.Tick(time.Duration(reportInterval) * time.Second)

	go services.UpdateStoragePeriodically(updateTicker, storage)
	go services.SendMetricsPeriodically(flagSendAddr, sendTicker, storage)
	fmt.Println("Agent started")
	// Ждем завершения работы программы (например, через сигнал ОС)
	select {}
}
