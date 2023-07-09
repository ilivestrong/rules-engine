// Code generated by mockery v2.18.0. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/ilivestrong/rules-engine/models"
	mock "github.com/stretchr/testify/mock"
)

// ApprovalRule is an autogenerated mock type for the ApprovalRule type
type ApprovalRule struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, applicant
func (_m *ApprovalRule) Execute(ctx context.Context, applicant models.Applicant) bool {
	ret := _m.Called(ctx, applicant)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, models.Applicant) bool); ok {
		r0 = rf(ctx, applicant)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockConstructorTestingTNewApprovalRule interface {
	mock.TestingT
	Cleanup(func())
}

// NewApprovalRule creates a new instance of ApprovalRule. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewApprovalRule(t mockConstructorTestingTNewApprovalRule) *ApprovalRule {
	mock := &ApprovalRule{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}