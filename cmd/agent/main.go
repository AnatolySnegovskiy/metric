package main

import (
	"context"
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/services/agent"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func handleError(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

func handleShutdownSignal(quit chan os.Signal) {
	<-quit
	fmt.Println("Agent stopped")
	os.Exit(0)
}

func main() {
	s := storages.NewMemStorage()
	s.AddMetric("gauge", metrics.NewGauge())
	s.AddMetric("counter", metrics.NewCounter())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go handleShutdownSignal(quit)

	fmt.Println("Agent started")
	c, err := NewConfig()
	handleError(err)

	handleError(
		agent.New(
			agent.Options{
				Storage:        s,
				PollInterval:   c.pollInterval,
				ReportInterval: c.reportInterval,
				SendAddr:       c.flagSendAddr,
			},
		).Run(context.Background()))
}
