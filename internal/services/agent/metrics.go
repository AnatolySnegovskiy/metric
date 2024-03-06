package agent

import (
	"context"
	"fmt"
	"net/http"
)

func sendMetricsPeriodically(ctx context.Context, addr string, s Storage) error {
	for storageType, storage := range s.GetList() {
		for metricName, metric := range storage.GetList() {
			err := sendMetric(ctx, addr, storageType, metricName, metric)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func sendMetric(ctx context.Context, addr string, storageType string, name string, metric any) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%v", addr, storageType, name, metric)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if err := resp.Body.Close(); err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
