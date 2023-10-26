// Code generated by mockery v2.35.2. DO NOT EDIT.

package tests

import (
	context "context"

	serviceaccounts "github.com/grafana/grafana/pkg/services/serviceaccounts"
	mock "github.com/stretchr/testify/mock"
)

// MockExtSvcAccountsService is an autogenerated mock type for the ExtSvcAccountsService type
type MockExtSvcAccountsService struct {
	mock.Mock
}

// ManageExtSvcAccount provides a mock function with given fields: ctx, cmd
func (_m *MockExtSvcAccountsService) ManageExtSvcAccount(ctx context.Context, cmd *serviceaccounts.ManageExtSvcAccountCmd) (int64, error) {
	ret := _m.Called(ctx, cmd)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *serviceaccounts.ManageExtSvcAccountCmd) (int64, error)); ok {
		return rf(ctx, cmd)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *serviceaccounts.ManageExtSvcAccountCmd) int64); ok {
		r0 = rf(ctx, cmd)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *serviceaccounts.ManageExtSvcAccountCmd) error); ok {
		r1 = rf(ctx, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RetrieveExtSvcAccount provides a mock function with given fields: ctx, orgID, saID
func (_m *MockExtSvcAccountsService) RetrieveExtSvcAccount(ctx context.Context, orgID int64, saID int64) (*serviceaccounts.ExtSvcAccount, error) {
	ret := _m.Called(ctx, orgID, saID)

	var r0 *serviceaccounts.ExtSvcAccount
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) (*serviceaccounts.ExtSvcAccount, error)); ok {
		return rf(ctx, orgID, saID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) *serviceaccounts.ExtSvcAccount); ok {
		r0 = rf(ctx, orgID, saID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*serviceaccounts.ExtSvcAccount)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, orgID, saID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockExtSvcAccountsService creates a new instance of MockExtSvcAccountsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockExtSvcAccountsService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockExtSvcAccountsService {
	mock := &MockExtSvcAccountsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
