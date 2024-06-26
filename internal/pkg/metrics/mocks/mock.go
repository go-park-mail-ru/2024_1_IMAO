// Code generated by MockGen. DO NOT EDIT.
// Source: mock_gen.go

// Package mock_metrics is a generated GoMock package.
package mock_metrics

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockDBMetrics is a mock of DBMetrics interface.
type MockDBMetrics struct {
	ctrl     *gomock.Controller
	recorder *MockDBMetricsMockRecorder
}

// MockDBMetricsMockRecorder is the mock recorder for MockDBMetrics.
type MockDBMetricsMockRecorder struct {
	mock *MockDBMetrics
}

// NewMockDBMetrics creates a new mock instance.
func NewMockDBMetrics(ctrl *gomock.Controller) *MockDBMetrics {
	mock := &MockDBMetrics{ctrl: ctrl}
	mock.recorder = &MockDBMetricsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBMetrics) EXPECT() *MockDBMetricsMockRecorder {
	return m.recorder
}

// AddDuration mocks base method.
func (m *MockDBMetrics) AddDuration(funcName string, duration time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddDuration", funcName, duration)
}

// AddDuration indicates an expected call of AddDuration.
func (mr *MockDBMetricsMockRecorder) AddDuration(funcName, duration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDuration", reflect.TypeOf((*MockDBMetrics)(nil).AddDuration), funcName, duration)
}

// IncreaseErrors mocks base method.
func (m *MockDBMetrics) IncreaseErrors(funcName string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IncreaseErrors", funcName)
}

// IncreaseErrors indicates an expected call of IncreaseErrors.
func (mr *MockDBMetricsMockRecorder) IncreaseErrors(funcName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncreaseErrors", reflect.TypeOf((*MockDBMetrics)(nil).IncreaseErrors), funcName)
}
