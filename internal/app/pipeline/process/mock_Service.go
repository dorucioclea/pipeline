// Code generated by mockery v1.0.0. DO NOT EDIT.

package process

import (
	context "context"

	pipeline "github.com/banzaicloud/pipeline/.gen/pipeline/pipeline"
	mock "github.com/stretchr/testify/mock"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

// CancelProcess provides a mock function with given fields: ctx, id
func (_m *MockService) CancelProcess(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetProcess provides a mock function with given fields: ctx, id
func (_m *MockService) GetProcess(ctx context.Context, id string) (pipeline.Process, error) {
	ret := _m.Called(ctx, id)

	var r0 pipeline.Process
	if rf, ok := ret.Get(0).(func(context.Context, string) pipeline.Process); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(pipeline.Process)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListProcesses provides a mock function with given fields: ctx, query
func (_m *MockService) ListProcesses(ctx context.Context, query pipeline.Process) ([]pipeline.Process, error) {
	ret := _m.Called(ctx, query)

	var r0 []pipeline.Process
	if rf, ok := ret.Get(0).(func(context.Context, pipeline.Process) []pipeline.Process); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]pipeline.Process)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pipeline.Process) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LogProcess provides a mock function with given fields: ctx, proc
func (_m *MockService) LogProcess(ctx context.Context, proc pipeline.Process) (pipeline.Process, error) {
	ret := _m.Called(ctx, proc)

	var r0 pipeline.Process
	if rf, ok := ret.Get(0).(func(context.Context, pipeline.Process) pipeline.Process); ok {
		r0 = rf(ctx, proc)
	} else {
		r0 = ret.Get(0).(pipeline.Process)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pipeline.Process) error); ok {
		r1 = rf(ctx, proc)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LogProcessEvent provides a mock function with given fields: ctx, proc
func (_m *MockService) LogProcessEvent(ctx context.Context, proc pipeline.ProcessEvent) (pipeline.ProcessEvent, error) {
	ret := _m.Called(ctx, proc)

	var r0 pipeline.ProcessEvent
	if rf, ok := ret.Get(0).(func(context.Context, pipeline.ProcessEvent) pipeline.ProcessEvent); ok {
		r0 = rf(ctx, proc)
	} else {
		r0 = ret.Get(0).(pipeline.ProcessEvent)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pipeline.ProcessEvent) error); ok {
		r1 = rf(ctx, proc)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
