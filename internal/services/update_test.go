// services_test.go
package services

import (
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetMetricType(metricType string) (storages.EntityMetric, error) {
	args := m.Called(metricType)
	return args.Get(0).(storages.EntityMetric), args.Error(1)
}

func (m *MockStorage) Process(key, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockStorage) GetList() map[string]storages.EntityMetric {
	args := m.Called()
	return args.Get(0).(map[string]storages.EntityMetric)

}

func (m *MockStorage) AddMetric(metricType string, metric storages.EntityMetric) {
}

type MockEntityMetric struct {
	mock.Mock
}

func (m *MockEntityMetric) Process(key, value string) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockEntityMetric) GetList() map[string]float64 {
	args := m.Called()
	return args.Get(0).(map[string]float64)
}

func TestUpdateStoragePeriodically(t *testing.T) {
	storage := &MockStorage{}
	entityMetric := &MockEntityMetric{}

	storage.On("GetMetricType", "gauge").Return(entityMetric, nil).Once()
	storage.On("GetMetricType", "counter").Return(entityMetric, nil).Once()
	entityMetric.On("Process", "PollCount", "1").Return(nil).Once()
	entityMetric.On("Process", "RandomValue", mock.Anything).Return(nil).Once()
	entityMetric.On("Process", mock.Anything, mock.Anything).Return(nil).Times(len(runtimeEntityArray))

	updateTicker := time.Tick(100 * time.Millisecond)

	// запускаем горутину
	go UpdateStoragePeriodically(updateTicker, storage)

	// ждем некоторое время для выполнения функции
	time.Sleep(150 * time.Millisecond)

	// проверяем, что все ожидаемые вызовы методов были выполнены
	storage.AssertExpectations(t)
}
