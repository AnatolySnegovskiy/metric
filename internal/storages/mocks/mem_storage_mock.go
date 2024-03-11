// Code generated by MockGen. DO NOT EDIT.
// Source: storages/mem_storage.go
//
// Generated by this command:
//
//	mockgen -source=storages/mem_storage.go -destination=storages/mocks/mem_storage_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	"github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockEntityMetric is a mock of EntityMetric interface.
type MockEntityMetric struct {
	ctrl     *gomock.Controller
	recorder *MockEntityMetricMockRecorder
}

// MockEntityMetricMockRecorder is the mock recorder for MockEntityMetric.
type MockEntityMetricMockRecorder struct {
	mock *MockEntityMetric
}

// NewMockEntityMetric creates a new mock instance.
func NewMockEntityMetric(ctrl *gomock.Controller) *MockEntityMetric {
	mock := &MockEntityMetric{ctrl: ctrl}
	mock.recorder = &MockEntityMetricMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEntityMetric) EXPECT() *MockEntityMetricMockRecorder {
	return m.recorder
}

// GetList mocks base method.
func (m *MockEntityMetric) GetList() map[string]float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList")
	ret0, _ := ret[0].(map[string]float64)
	return ret0
}

// GetList indicates an expected call of GetList.
func (mr *MockEntityMetricMockRecorder) GetList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockEntityMetric)(nil).GetList))
}

// Process mocks base method.
func (m *MockEntityMetric) Process(name, data string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", name, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Process indicates an expected call of Process.
func (mr *MockEntityMetricMockRecorder) Process(name, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockEntityMetric)(nil).Process), name, data)
}
