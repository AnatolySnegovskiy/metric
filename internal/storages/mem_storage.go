package storages

import (
	"errors"
)

type EntityMetric interface {
	Process(name string, data string) error
	GetList() (map[string]float64, error)
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

func (m *MemStorage) GetList() map[string]EntityMetric {
	return m.metrics
}
