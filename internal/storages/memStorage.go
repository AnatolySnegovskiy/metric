package storages

import (
	"errors"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"log"
)

type EntityMetric interface {
	Process(name string, data string) error
}

type MemStorage struct {
	metrics map[string]EntityMetric
}

func NewMemStorage() *MemStorage {
	storage := &MemStorage{
		metrics: make(map[string]EntityMetric),
	}
	storage.metrics["gauge"] = metrics.NewGauge()
	storage.metrics["counter"] = metrics.NewCounter()
	return storage
}

func (m *MemStorage) GetMetricType(metricType string) (EntityMetric, error) {
	if m.metrics[metricType] == nil {
		return nil, errors.New("metric type not found")
	}

	return m.metrics[metricType], nil
}

func (m *MemStorage) Log() {
	for _, v := range m.metrics {
		log.Println(v)
	}
}
