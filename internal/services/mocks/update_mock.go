// Code generated by MockGen. DO NOT EDIT.
// Source: update.go
//
// Generated by this command:
//
//	mockgen -source=update.go -destination=mocks/update_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	storages "github.com/AnatolySnegovskiy/metric/internal/storages"
	gomock "go.uber.org/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddMetric mocks base method.
func (m *MockStorage) AddMetric(metricType string, metric storages.EntityMetric) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddMetric", metricType, metric)
}

// AddMetric indicates an expected call of AddMetric.
func (mr *MockStorageMockRecorder) AddMetric(metricType, metric any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMetric", reflect.TypeOf((*MockStorage)(nil).AddMetric), metricType, metric)
}

// GetList mocks base method.
func (m *MockStorage) GetList() map[string]storages.EntityMetric {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList")
	ret0, _ := ret[0].(map[string]storages.EntityMetric)
	return ret0
}

// GetList indicates an expected call of GetList.
func (mr *MockStorageMockRecorder) GetList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockStorage)(nil).GetList))
}

// GetMetricType mocks base method.
func (m *MockStorage) GetMetricType(metricType string) (storages.EntityMetric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetricType", metricType)
	ret0, _ := ret[0].(storages.EntityMetric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetricType indicates an expected call of GetMetricType.
func (mr *MockStorageMockRecorder) GetMetricType(metricType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetricType", reflect.TypeOf((*MockStorage)(nil).GetMetricType), metricType)
}
