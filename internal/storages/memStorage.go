package storages

import (
	"errors"
	"log"
	"strconv"
)

type MemStorage struct {
	metrics map[string]StorageInterface
}

func NewMemStorage() *MemStorage {
	storage := &MemStorage{
		metrics: make(map[string]StorageInterface),
	}
	storage.metrics["gauge"] = &gauge{list: make(map[string]float64)}
	storage.metrics["counter"] = &counter{list: make(map[string]int)}
	return storage
}

func (m *MemStorage) GetMetricType(metricType string) (StorageInterface, error) {
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

type StorageInterface interface {
	Process(name string, data string) error
}

type gauge struct {
	list map[string]float64
}

func (g *gauge) Process(name string, data string) error {
	floatValue, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return errors.New("metric value is not float64")
	}

	g.list[name] = floatValue
	return nil
}

type counter struct {
	list map[string]int
}

func (c *counter) Process(name string, data string) error {
	intValue, err := strconv.Atoi(data)
	if err != nil {
		return errors.New("metric value is not int")
	}

	c.list[name] += intValue
	return nil
}
