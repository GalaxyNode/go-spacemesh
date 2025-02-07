// Code generated by MockGen. DO NOT EDIT.
// Source: ./interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	types "github.com/spacemeshos/go-spacemesh/common/types"
)

// MockmeshProvider is a mock of meshProvider interface.
type MockmeshProvider struct {
	ctrl     *gomock.Controller
	recorder *MockmeshProviderMockRecorder
}

// MockmeshProviderMockRecorder is the mock recorder for MockmeshProvider.
type MockmeshProviderMockRecorder struct {
	mock *MockmeshProvider
}

// NewMockmeshProvider creates a new mock instance.
func NewMockmeshProvider(ctrl *gomock.Controller) *MockmeshProvider {
	mock := &MockmeshProvider{ctrl: ctrl}
	mock.recorder = &MockmeshProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmeshProvider) EXPECT() *MockmeshProviderMockRecorder {
	return m.recorder
}

// AddTXsFromProposal mocks base method.
func (m *MockmeshProvider) AddTXsFromProposal(arg0 context.Context, arg1 types.LayerID, arg2 types.ProposalID, arg3 []types.TransactionID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTXsFromProposal", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTXsFromProposal indicates an expected call of AddTXsFromProposal.
func (mr *MockmeshProviderMockRecorder) AddTXsFromProposal(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTXsFromProposal", reflect.TypeOf((*MockmeshProvider)(nil).AddTXsFromProposal), arg0, arg1, arg2, arg3)
}

// MockeligibilityValidator is a mock of eligibilityValidator interface.
type MockeligibilityValidator struct {
	ctrl     *gomock.Controller
	recorder *MockeligibilityValidatorMockRecorder
}

// MockeligibilityValidatorMockRecorder is the mock recorder for MockeligibilityValidator.
type MockeligibilityValidatorMockRecorder struct {
	mock *MockeligibilityValidator
}

// NewMockeligibilityValidator creates a new mock instance.
func NewMockeligibilityValidator(ctrl *gomock.Controller) *MockeligibilityValidator {
	mock := &MockeligibilityValidator{ctrl: ctrl}
	mock.recorder = &MockeligibilityValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockeligibilityValidator) EXPECT() *MockeligibilityValidatorMockRecorder {
	return m.recorder
}

// CheckEligibility mocks base method.
func (m *MockeligibilityValidator) CheckEligibility(arg0 context.Context, arg1 *types.Ballot) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckEligibility", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckEligibility indicates an expected call of CheckEligibility.
func (mr *MockeligibilityValidatorMockRecorder) CheckEligibility(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckEligibility", reflect.TypeOf((*MockeligibilityValidator)(nil).CheckEligibility), arg0, arg1)
}
