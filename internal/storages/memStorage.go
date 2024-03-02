package storages

import (
	"errors"
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
	return storage
}

func (m *MemStorage) AddMetric(metricType string, metric EntityMetric) {
	m.metrics[metricType] = metric
}

func (m *MemStorage) GetMetricType(metricType string) (EntityMetric, error) {
	mt, ok := m.metrics[metricType]
	if !ok {
		return nil, errors.New("metric type not found")
	}

	return mt, nil
}

func (m *MemStorage) Log() {
	for _, v := range m.metrics {
		log.Println(v)
	}
}
