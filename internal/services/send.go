package services

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"time"
)

func SendMetricsPeriodically(addr string, ticker <-chan time.Time, s *storages.MemStorage) {
	for range ticker {
		for storageType, storage := range s.GetList() {
			for metricName, metric := range storage.GetList() {
				err := sendMetric(addr, storageType, metricName, metric)
				if err != nil {
					fmt.Println("Ошибка отправки метрик:", err)
				}
			}
		}

		fmt.Println("Метрика отправлена")
	}
}
