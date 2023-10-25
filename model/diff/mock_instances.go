// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aserto-dev/azm/model/diff (interfaces: Instances)

// Package diff is a generated GoMock package.
package diff

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockInstances is a mock of Instances interface
type MockInstances struct {
	ctrl     *gomock.Controller
	recorder *MockInstancesMockRecorder
}

// MockInstancesMockRecorder is the mock recorder for MockInstances
type MockInstancesMockRecorder struct {
	mock *MockInstances
}

// NewMockInstances creates a new mock instance
func NewMockInstances(ctrl *gomock.Controller) *MockInstances {
	mock := &MockInstances{ctrl: ctrl}
	mock.recorder = &MockInstancesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInstances) EXPECT() *MockInstancesMockRecorder {
	return m.recorder
}

// ObjectTypes mocks base method
func (m *MockInstances) ObjectTypes() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ObjectTypes")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ObjectTypes indicates an expected call of ObjectTypes
func (mr *MockInstancesMockRecorder) ObjectTypes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ObjectTypes", reflect.TypeOf((*MockInstances)(nil).ObjectTypes))
}

// RelationTypes mocks base method
func (m *MockInstances) RelationTypes() ([]*RelationKind, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RelationTypes")
	ret0, _ := ret[0].([]*RelationKind)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RelationTypes indicates an expected call of RelationTypes
func (mr *MockInstancesMockRecorder) RelationTypes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RelationTypes", reflect.TypeOf((*MockInstances)(nil).RelationTypes))
}
