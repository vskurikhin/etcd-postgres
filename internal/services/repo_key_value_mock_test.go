// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/domain/repo.go
//
// Generated by this command:
//
//	mockgen -source=./internal/domain/repo.go -package=services
//

// Package services is a generated GoMock package.
package services

import (
	context "context"
	reflect "reflect"

	domain "github.com/victor-skurikhin/etcd-client/v1/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockRepo is a mock of Repo interface.
type MockRepo[A domain.Actioner[T, U], T domain.Ptr[U], U domain.Entity] struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder[A, T, U]
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder[A domain.Actioner[T, U], T domain.Ptr[U], U domain.Entity] struct {
	mock *MockRepo[A, T, U]
}

// NewMockRepo creates a new mock instance.
func NewMockRepo[A domain.Actioner[T, U], T domain.Ptr[U], U domain.Entity](ctrl *gomock.Controller) *MockRepo[A, T, U] {
	mock := &MockRepo[A, T, U]{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder[A, T, U]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo[A, T, U]) EXPECT() *MockRepoMockRecorder[A, T, U] {
	return m.recorder
}

// Do mocks base method.
func (m *MockRepo[A, T, U]) Do(ctx context.Context, action A, unit U, scan func(domain.Scanner) U) (U, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", ctx, action, unit, scan)
	ret0, _ := ret[0].(U)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockRepoMockRecorder[A, T, U]) Do(ctx, action, unit, scan any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockRepo[A, T, U])(nil).Do), ctx, action, unit, scan)
}

// Get mocks base method.
func (m *MockRepo[A, T, U]) Get(ctx context.Context, action A, unit U, scan func(domain.Scanner) U) ([]U, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, action, unit, scan)
	ret0, _ := ret[0].([]U)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRepoMockRecorder[A, T, U]) Get(ctx, action, unit, scan any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepo[A, T, U])(nil).Get), ctx, action, unit, scan)
}
