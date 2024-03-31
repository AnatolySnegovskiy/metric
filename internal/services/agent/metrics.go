package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/services/dto"
	"github.com/mailru/easyjson"
	"net/http"
)

func (a *Agent) sendMetricsPeriodically(ctx context.Context) error {
	metricDto := &dto.Metrics{}
	for storageType, storage := range a.storage.GetList() {
		if storage == nil {
			continue
		}

		for metricName, metric := range storage.GetList() {
			metricDto.MType = storageType
			metricDto.ID = metricName

			if storageType == "counter" {
				iv := int64(metric)
				metricDto.Delta = &iv
			} else {
				metricDto.Value = &metric
			}

			var buf bytes.Buffer
			gw := gzip.NewWriter(&buf)
			defer gw.Close()

			body, _ := easyjson.Marshal(metricDto)
			if _, err := gw.Write(body); err != nil {
				return err
			}

			if err := gw.Close(); err != nil {
				return err
			}

			url := fmt.Sprintf("http://%s/update/", a.sendAddr)
			req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
			req.Header.Set("Content-Encoding", "gzip")
			req.Header.Set("Content-Type", "application/json")
			resp, _ := a.client.Do(req)

			if resp != nil {
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
				}
			}
		}
	}

	return nil
}
