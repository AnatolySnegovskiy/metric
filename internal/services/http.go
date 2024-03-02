package services

import (
	"fmt"
	"io"
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

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	return nil
}
