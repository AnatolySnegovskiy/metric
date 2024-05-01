package agent

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (a *Agent) Run(ctx context.Context) error {
	pollTicker := time.NewTicker(time.Duration(a.pollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(a.reportInterval) * time.Second)
	retrievableCounter := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-pollTicker.C:
			err := a.updateStoragePeriodically(ctx)
			if err != nil {
				return fmt.Errorf("error occurred while updating storage: %w", err)
			}
			log.Println("storage updated")
		case <-pollTicker.C:
			err := a.updateGopsutil(ctx)
			if err != nil {
				return fmt.Errorf("error occurred while updating storage: %w", err)
			}
			log.Println("storage updated")
		case <-reportTicker.C:
			err := a.sendMetricsPeriodically(ctx)
			if err != nil {
				if retrievableCounter < a.maxRetries {
					retrievableCounter++
					log.Println(err)
					continue
				}
				return fmt.Errorf("error occurred while sending metrics: %w", err)
			}
			retrievableCounter = 0
			log.Println("metrics sent")
		}
	}
}
