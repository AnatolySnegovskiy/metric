package storages_test

import (
	"testing"

	"github.com/AnatolySnegovskiy/metric/internal/storages"
)

type MockMetric struct{}

func (m *MockMetric) Process(name, data string) error {
	// Implement mock behavior for Process
	return nil
}

func (m *MockMetric) GetList() map[string]float64 {
	// Implement mock behavior for GetList
	return nil
}
func TestMemStorage_GetMetricType(t *testing.T) {
	storage := storages.NewMemStorage()
	mockMetric := &MockMetric{}
	storage.AddMetric("test", mockMetric)

	// Test case 1: metric type exists
	_, err := storage.GetMetricType("test")
	if err != nil {
		t.Errorf("expected metric type to exist, got error: %v", err)
	}

	// Test case 2: metric type does not exist
	_, err = storage.GetMetricType("invalid")
	if err == nil {
		t.Error("expected metric type not to exist, but got no error")
	}
}
