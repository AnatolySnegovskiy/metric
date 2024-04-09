package storages_test

import (
	"github.com/AnatolySnegovskiy/metric/internal/mocks"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestMemStorage(t *testing.T) {

	tests := []struct {
		name     string
		setup    func(storage *storages.MemStorage, ctrl *gomock.Controller)
		validate func(t *testing.T, storage *storages.MemStorage, ctrl *gomock.Controller)
	}{
		{
			name: "AddMetric and GetMetricType",
			setup: func(storage *storages.MemStorage, ctrl *gomock.Controller) {
				mockMetric := mocks.NewMockEntityMetric(ctrl)
				mockMetricType := "mockType"
				storage.AddMetric(mockMetricType, mockMetric)
			},
			validate: func(t *testing.T, storage *storages.MemStorage, ctrl *gomock.Controller) {
				retrievedMetric, err := storage.GetMetricType("mockType")
				assert.Nil(t, err, "Unexpected error")
				assert.NotNil(t, retrievedMetric, "Retrieved metric is nil")
			},
		},
		{
			name:  "GetMetricType_NotFound",
			setup: func(storage *storages.MemStorage, ctrl *gomock.Controller) {},
			validate: func(t *testing.T, storage *storages.MemStorage, ctrl *gomock.Controller) {
				_, err := storage.GetMetricType("nonExistentType")
				assert.NotNil(t, err, "Expected error for non-existent metric type")
				assert.ErrorContains(t, err, "metric type not found")
			},
		},
		{
			name: "GetList",
			setup: func(storage *storages.MemStorage, ctrl *gomock.Controller) {
				mockMetric1 := mocks.NewMockEntityMetric(ctrl)
				mockMetric2 := mocks.NewMockEntityMetric(ctrl)
				storage.AddMetric("type1", mockMetric1)
				storage.AddMetric("type2", mockMetric2)
			},
			validate: func(t *testing.T, storage *storages.MemStorage, ctrl *gomock.Controller) {
				metricList := storage.GetList()
				assert.Len(t, metricList, 2, "Expected 2 metrics in the list")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := storages.NewMemStorage()
			tt.setup(storage, ctrl)
			tt.validate(t, storage, ctrl)
		})
	}
}
