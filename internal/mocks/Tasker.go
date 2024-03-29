// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	entities "github.com/Amirhossein2000/RequestTasker/internal/domain/entities"
	mock "github.com/stretchr/testify/mock"
)

// TaskerMock is an autogenerated mock type for the Tasker type
type TaskerMock struct {
	mock.Mock
}

// Process provides a mock function with given fields: ctx, task
func (_m *TaskerMock) Process(ctx context.Context, task entities.Task) error {
	ret := _m.Called(ctx, task)

	if len(ret) == 0 {
		panic("no return value specified for Process")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entities.Task) error); ok {
		r0 = rf(ctx, task)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTaskerMock creates a new instance of TaskerMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTaskerMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *TaskerMock {
	mock := &TaskerMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
