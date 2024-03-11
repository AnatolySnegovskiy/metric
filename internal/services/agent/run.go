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

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-pollTicker.C:
			err := a.updateStoragePeriodically()
			if err != nil {
				return fmt.Errorf("error occurred while updating storage: %w", err)
			}
			log.Println("storage updated")
		case <-reportTicker.C:
			err := a.sendMetricsPeriodically(ctx)
			if err != nil {
				return fmt.Errorf("error occurred while sending metrics: %w", err)
			}
			log.Println("metrics sent")
		}
	}
}
