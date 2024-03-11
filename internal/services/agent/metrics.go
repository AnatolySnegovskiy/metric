package agent

import (
	"context"
	"fmt"
	"net/http"
)

func (a *Agent) sendMetricsPeriodically(ctx context.Context) error {
	for storageType, storage := range a.storage.GetList() {
		for metricName, metric := range storage.GetList() {
			url := fmt.Sprintf("http://%s/update/%s/%s/%v", a.sendAddr, storageType, metricName, metric)
			req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
			resp, err := a.client.Do(req)

			if err != nil {
				return err
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			}

			return nil
		}
	}

	return fmt.Errorf("no metrics to send")
}
