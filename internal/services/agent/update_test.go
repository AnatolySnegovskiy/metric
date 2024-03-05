package agent

import (
	mocks2 "github.com/AnatolySnegovskiy/metric/internal/services/agent/mocks"
	"github.com/AnatolySnegovskiy/metric/internal/storages/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUpdateStoragePeriodically(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks2.NewMockStorage(ctrl)
	mockGauge := mocks.NewMockEntityMetric(ctrl)
	mockCounter := mocks.NewMockEntityMetric(ctrl)

	mockStorage.EXPECT().GetMetricType("gauge").Return(mockGauge, nil)
	mockStorage.EXPECT().GetMetricType("counter").Return(mockCounter, nil)

	mockCounter.EXPECT().Process("PollCount", "1").Return(nil)
	mockGauge.EXPECT().Process("RandomValue", gomock.Any()).Return(nil).AnyTimes()

	for _, field := range runtimeEntityArray {
		mockGauge.EXPECT().Process(field, gomock.Any()).Return(nil)
	}

	err := UpdateStoragePeriodically(mockStorage)
	assert.NoError(t, err, "Expected no error")
}
