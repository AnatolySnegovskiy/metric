package agent

import (
	"context"
	"errors"
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/mocks"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	options := Options{
		Storage:        mocks.NewMockStorage(gomock.NewController(t)),
		PollInterval:   10,
		ReportInterval: 20,
		SendAddr:       "example.com:1234",
	}

	agent := New(options)
	assert.NotNil(t, agent.storage, "storage should not be nil")
	assert.NotNil(t, agent.pollInterval, "pollInterval should not be nil")
	assert.NotNil(t, agent.reportInterval, "reportInterval should not be nil")
	assert.Equal(t, "example.com:1234", agent.sendAddr, "send address should be example.com:1234")
}

func TestAgent(t *testing.T) {
	testCases := []struct {
		name          string
		statusCode    int
		doReturnError error
		expectedErr   bool
		mockStorage   func() *storages.MemStorage
	}{
		{"SuccessNil", http.StatusOK, nil, false, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			mockStorage.AddMetric("gauge", metrics.NewGauge(nil))
			mockStorage.AddMetric("counter", metrics.NewCounter(nil))
			mockStorage.AddMetric("nil", nil)
			return mockStorage
		}},
		{"Success", http.StatusOK, nil, false, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			mockStorage.AddMetric("gauge", metrics.NewGauge(nil))
			mockStorage.AddMetric("counter", metrics.NewCounter(nil))
			return mockStorage
		}},
		{"ErrorPoll", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			return mockStorage
		}},
		{"ErrorPollCounter", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			mockStorage.AddMetric("gauge", metrics.NewGauge(nil))
			return mockStorage
		}},
		{"ErrorPollGauge", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			mockStorage.AddMetric("counter", metrics.NewCounter(nil))
			return mockStorage
		}},

		{"ErrorPollPollCount", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			ctrl := gomock.NewController(t)
			mockEntity := mocks.NewMockEntityMetric(ctrl)
			mockEntity.EXPECT().Process(gomock.Any(), "PollCount", gomock.Any()).Return(
				errors.New("some error"),
			).AnyTimes().MinTimes(1)
			mockEntity.EXPECT().Process(gomock.Any(), gomock.Not("PollCount"), gomock.Any()).Return(
				nil,
			).AnyTimes()
			mockEntity.EXPECT().GetList(gomock.Any()).Return(map[string]float64{}, nil).AnyTimes()
			mockStorage.AddMetric("counter", mockEntity)
			mockStorage.AddMetric("gauge", metrics.NewGauge(nil))
			return mockStorage
		}},
		{"ErrorTotalMemory", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			ctrl := gomock.NewController(t)
			mockEntity := mocks.NewMockEntityMetric(ctrl)
			mockEntity.EXPECT().Process(gomock.Any(), "TotalMemory", gomock.Any()).Return(
				errors.New("some error"),
			).AnyTimes().MinTimes(1)
			mockEntity.EXPECT().Process(gomock.Any(), gomock.Not("TotalMemory"), gomock.Any()).Return(
				nil,
			).AnyTimes()
			mockEntity.EXPECT().GetList(gomock.Any()).Return(map[string]float64{}, nil).AnyTimes()
			mockStorage.AddMetric("counter", metrics.NewGauge(nil))
			mockStorage.AddMetric("gauge", mockEntity)
			return mockStorage
		}},
		{"ErrorFreeMemory", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			ctrl := gomock.NewController(t)
			mockEntity := mocks.NewMockEntityMetric(ctrl)
			mockEntity.EXPECT().Process(gomock.Any(), "FreeMemory", gomock.Any()).Return(
				errors.New("some error"),
			).AnyTimes().MinTimes(1)
			mockEntity.EXPECT().Process(gomock.Any(), gomock.Not("FreeMemory"), gomock.Any()).Return(
				nil,
			).AnyTimes()
			mockEntity.EXPECT().GetList(gomock.Any()).Return(map[string]float64{}, nil).AnyTimes()

			mockStorage.AddMetric("counter", metrics.NewGauge(nil))
			mockStorage.AddMetric("gauge", mockEntity)
			return mockStorage
		}},
		{"ErrorCPUutilization1", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			ctrl := gomock.NewController(t)
			mockEntity := mocks.NewMockEntityMetric(ctrl)
			mockEntity.EXPECT().Process(gomock.Any(), gomock.Not("CPUutilization1"), gomock.Any()).Return(
				nil,
			).AnyTimes()
			mockEntity.EXPECT().Process(gomock.Any(), "CPUutilization1", gomock.Any()).Return(
				errors.New("some error"),
			).AnyTimes().MinTimes(1)
			mockEntity.EXPECT().GetList(gomock.Any()).Return(map[string]float64{}, nil).AnyTimes()
			mockStorage.AddMetric("counter", metrics.NewGauge(nil))
			mockStorage.AddMetric("gauge", mockEntity)
			return mockStorage
		}},
		{"ErrorPollRandomValue", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			ctrl := gomock.NewController(t)
			mockEntity := mocks.NewMockEntityMetric(ctrl)
			mockEntity.EXPECT().Process(gomock.Any(), "RandomValue", gomock.Any()).Return(
				errors.New("some error"),
			).AnyTimes().MinTimes(1)
			mockEntity.EXPECT().Process(gomock.Any(), gomock.Not("RandomValue"), gomock.Any()).Return(
				nil,
			).AnyTimes()
			mockEntity.EXPECT().GetList(gomock.Any()).Return(map[string]float64{}, nil).AnyTimes()
			mockStorage.AddMetric("counter", metrics.NewGauge(nil))
			mockStorage.AddMetric("gauge", mockEntity)
			return mockStorage
		}},
		{"ErrorPollRuntimeEntityArray", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			ctrl := gomock.NewController(t)
			mockEntity := mocks.NewMockEntityMetric(ctrl)
			mockEntity.EXPECT().Process(gomock.Any(), "Alloc", gomock.Any()).Return(
				errors.New("some error"),
			).AnyTimes().MinTimes(1)
			mockEntity.EXPECT().Process(gomock.Any(), gomock.Not("Alloc"), gomock.Any()).Return(
				nil,
			).AnyTimes()
			mockEntity.EXPECT().GetList(gomock.Any()).Return(map[string]float64{}, nil).AnyTimes()
			mockStorage.AddMetric("counter", metrics.NewGauge(nil))
			mockStorage.AddMetric("gauge", mockEntity)
			return mockStorage
		}},
		{"ErrorReport", http.StatusBadRequest, fmt.Errorf("some error"), true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			mockStorage.AddMetric("gauge", metrics.NewGauge(nil))
			mockStorage.AddMetric("counter", metrics.NewCounter(nil))
			return mockStorage
		}},
		{"StatusBadRequest", http.StatusBadRequest, nil, true, func() *storages.MemStorage {
			mockStorage := storages.NewMemStorage()
			mockStorage.AddMetric("gauge", metrics.NewGauge(nil))
			mockStorage.AddMetric("counter", metrics.NewCounter(nil))
			return mockStorage
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
			defer cancel()
			go func() {
				<-ctx.Done()

				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				httpClient := mocks.NewMockHTTPClient(ctrl)

				resp := &http.Response{StatusCode: tc.statusCode, Body: http.NoBody}
				httpClient.EXPECT().Do(gomock.Any()).Return(resp, tc.doReturnError).AnyTimes()

				a := &Agent{
					storage:        tc.mockStorage(),
					sendAddr:       "testAddr",
					client:         httpClient,
					pollInterval:   1,
					reportInterval: 1,
					maxRetries:     1,
					shaKey:         "testKey",
				}
				err := a.Run(context.Background())
				if tc.expectedErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}()
		})
	}
}

func TestAgentReportTickerEmpty(t *testing.T) {
	t.Run("EmptyStorage", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		httpClient := mocks.NewMockHTTPClient(ctrl)
		resp := &http.Response{}
		httpClient.EXPECT().Do(gomock.Any()).Return(resp, nil).AnyTimes()
		a := &Agent{
			storage:        storages.NewMemStorage(),
			sendAddr:       "testAddr",
			client:         httpClient,
			pollInterval:   1,
			reportInterval: 1,
			shaKey:         "testKey",
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := a.Run(ctx)

		assert.Error(t, err)
	})
}
