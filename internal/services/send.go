package services

import (
	"github.com/AnatolySnegovskiy/metric/internal/storages"
)

func SendMetricsPeriodically(addr string, s *storages.MemStorage) error {
	for storageType, storage := range s.GetList() {
		for metricName, metric := range storage.GetList() {
			err := sendMetric(addr, storageType, metricName, metric)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
