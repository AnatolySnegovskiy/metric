package agent

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (a *Agent) Run(ctx context.Context) error {
	updateTicker := time.Tick(time.Duration(a.pollInterval) * time.Second)
	sendTicker := time.Tick(time.Duration(a.reportInterval) * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-updateTicker:
			err := a.updateStoragePeriodically()
			if err != nil {
				return fmt.Errorf("error occurred while updating storage: %w", err)
			}
			log.Println("storage updated")
		case <-sendTicker:
			err := a.sendMetricsPeriodically(ctx)
			if err != nil {
				return fmt.Errorf("error occurred while sending metrics: %w", err)
			}
			log.Println("metrics sent")
		}
	}
}
