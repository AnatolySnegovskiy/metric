package main

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/services"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handleError(err error, message string) {
	if err != nil {
		log.Println(message + err.Error())
		os.Exit(1)
	}
}

func handleShutdownSignal(quit chan os.Signal) {
	<-quit
	fmt.Println("Agent stopped")
	os.Exit(0)
}

func main() {
	if err := parseFlags(); err != nil {
		handleError(err, "error occurred while parsing flags: ")
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go handleShutdownSignal(quit)

	s := storages.NewMemStorage()
	s.AddMetric("gauge", metrics.NewGauge())
	s.AddMetric("counter", metrics.NewCounter())
	fmt.Println("Agent started")

	runUpdateAndSendLoop(s)
}

func runUpdateAndSendLoop(s *storages.MemStorage) {
	updateTicker := time.Tick(time.Duration(pollInterval) * time.Second)
	sendTicker := time.Tick(time.Duration(reportInterval) * time.Second)

	for {
		select {
		case <-updateTicker:
			err := services.UpdateStoragePeriodically(s)
			if err != nil {
				handleError(err, "error occurred while updating storage: ")
			}
			log.Println("storage updated")
		case <-sendTicker:
			err := services.SendMetricsPeriodically(flagSendAddr, s)
			if err != nil {
				handleError(err, "error occurred while sending metrics: ")
			}
			log.Println("metrics sent")
		}
	}
}
