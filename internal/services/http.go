package services

import (
	"context"
	"fmt"
	"net/http"
)

func sendMetric(addr string, storageType string, name string, metric any) error {
	url := fmt.Sprintf("http://%s/update/%s/%s/%v", addr, storageType, name, metric)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, nil)
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
