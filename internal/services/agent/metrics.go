package agent

import (
	"bytes"
	"context"
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/services/dto"
	"github.com/mailru/easyjson"
	"net/http"
)

func (a *Agent) sendMetricsPeriodically(ctx context.Context) error {
	metricDto := &dto.Metrics{}
	for storageType, storage := range a.storage.GetList() {
		for metricName, metric := range storage.GetList() {
			metricDto.MType = storageType
			metricDto.ID = metricName

			if storageType == "counter" {
				iv := int64(metric)
				metricDto.Delta = &iv
			} else {
				metricDto.Value = &metric
			}

			url := fmt.Sprintf("http://%s/update/", a.sendAddr)
			body, _ := easyjson.Marshal(metricDto)
			req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := a.client.Do(req)
			
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			}
		}
	}

	return nil
}
