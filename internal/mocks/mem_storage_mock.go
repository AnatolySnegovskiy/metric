// Code generated by MockGen. DO NOT EDIT.
// Source: storages/mem_storage.go
//
// Generated by this command:
//
//	mockgen -source=storages/mem_storage.go -destination=mocks/mem_storage_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
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
func (m *MockEntityMetric) GetList(ctx context.Context) (map[string]float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList", ctx)
	ret0, _ := ret[0].(map[string]float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetList indicates an expected call of GetList.
func (mr *MockEntityMetricMockRecorder) GetList(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockEntityMetric)(nil).GetList), ctx)
}

// Process mocks base method.
func (m *MockEntityMetric) Process(ctx context.Context, name, data string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", ctx, name, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Process indicates an expected call of Process.
func (mr *MockEntityMetricMockRecorder) Process(ctx, name, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockEntityMetric)(nil).Process), ctx, name, data)
}

// ProcessMassive mocks base method.
func (m *MockEntityMetric) ProcessMassive(ctx context.Context, data map[string]float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessMassive", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessMassive indicates an expected call of ProcessMassive.
func (mr *MockEntityMetricMockRecorder) ProcessMassive(ctx, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessMassive", reflect.TypeOf((*MockEntityMetric)(nil).ProcessMassive), ctx, data)
}
