package services

import (
	"fmt"
	"net/http"
)

const serverAddress = "http://localhost:8080"

func sendMetric(storageType string, name string, metric any) error {
	url := fmt.Sprintf("%s/update/%s/%s/%v", serverAddress, storageType, name, metric)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
