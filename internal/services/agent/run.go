package agent

import (
	"context"
	"fmt"
	"log"
)

func (a *Agent) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-a.pollInterval:
			err := a.updateStoragePeriodically()
			if err != nil {
				return fmt.Errorf("error occurred while updating storage: %w", err)
			}
			log.Println("storage updated")
		case <-a.reportInterval:
			err := a.sendMetricsPeriodically(ctx)
			if err != nil {
				return fmt.Errorf("error occurred while sending metrics: %w", err)
			}
			log.Println("metrics sent")
		}
	}
}
