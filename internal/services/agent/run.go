package agent

import (
	"context"
	"fmt"
	"log"
)

func (a *Agent) Run(ctx context.Context) error {
	s := a.storage

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-a.pollInterval:
			err := updateStoragePeriodically(s)
			if err != nil {
				return fmt.Errorf("error occurred while updating storage: %w", err)
			}
			log.Println("storage updated")
		case <-a.reportInterval:
			err := sendMetricsPeriodically(ctx, a.flagSendAddr, s)
			if err != nil {
				return fmt.Errorf("error occurred while sending metrics: %w", err)
			}
			log.Println("metrics sent")
		}
	}
}
