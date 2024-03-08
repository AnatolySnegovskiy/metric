package agent

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (a *Agent) Run(ctx context.Context) error {

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(time.Duration(a.pollInterval) * time.Second):
			err := a.updateStoragePeriodically()
			if err != nil {
				return fmt.Errorf("error occurred while updating storage: %w", err)
			}
			log.Println("storage updated")
		case <-time.Tick(time.Duration(a.reportInterval) * time.Second):
			err := a.sendMetricsPeriodically(ctx)
			if err != nil {
				return fmt.Errorf("error occurred while sending metrics: %w", err)
			}
			log.Println("metrics sent")
		}
	}
}
