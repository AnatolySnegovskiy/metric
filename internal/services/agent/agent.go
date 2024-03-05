package agent

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Agent struct {
	storage Storage
}

func New(storage Storage) *Agent {
	return &Agent{
		storage: storage,
	}
}

func handleShutdownSignal(quit chan os.Signal) {
	<-quit
	fmt.Println("Agent stopped")
	os.Exit(0)
}

func (a *Agent) Run() error {
	err := parseFlags()

	if err != nil {
		return fmt.Errorf("error occurred while parsing flags: %w", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go handleShutdownSignal(quit)

	s := a.storage
	updateTicker := time.Tick(time.Duration(pollInterval) * time.Second)
	sendTicker := time.Tick(time.Duration(reportInterval) * time.Second)

	for {
		select {
		case <-updateTicker:
			err := UpdateStoragePeriodically(s)
			if err != nil {
				return fmt.Errorf("error occurred while updating storage: %w", err)
			}
			log.Println("storage updated")
		case <-sendTicker:
			err := SendMetricsPeriodically(flagSendAddr, s)
			if err != nil {
				return fmt.Errorf("error occurred while sending metrics: %w", err)
			}
			log.Println("metrics sent")
		}
	}
}
