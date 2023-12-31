// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package user

import (
	"context"
	"github.com/pollenjp/gameserver-go/api/entity"
	"sync"
)

// Ensure, that GetUserServiceMock does implement GetUserService.
// If this is not the case, regenerate this file with moq.
var _ GetUserService = &GetUserServiceMock{}

// GetUserServiceMock is a mock implementation of GetUserService.
//
//	func TestSomethingThatUsesGetUserService(t *testing.T) {
//
//		// make and configure a mocked GetUserService
//		mockedGetUserService := &GetUserServiceMock{
//			GetUserFunc: func(ctx context.Context, userId entity.UserId) (*entity.User, error) {
//				panic("mock out the GetUser method")
//			},
//		}
//
//		// use mockedGetUserService in code that requires GetUserService
//		// and then make assertions.
//
//	}
type GetUserServiceMock struct {
	// GetUserFunc mocks the GetUser method.
	GetUserFunc func(ctx context.Context, userId entity.UserId) (*entity.User, error)

	// calls tracks calls to the methods.
	calls struct {
		// GetUser holds details about calls to the GetUser method.
		GetUser []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserId is the userId argument value.
			UserId entity.UserId
		}
	}
	lockGetUser sync.RWMutex
}

// GetUser calls GetUserFunc.
func (mock *GetUserServiceMock) GetUser(ctx context.Context, userId entity.UserId) (*entity.User, error) {
	if mock.GetUserFunc == nil {
		panic("GetUserServiceMock.GetUserFunc: method is nil but GetUserService.GetUser was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		UserId entity.UserId
	}{
		Ctx:    ctx,
		UserId: userId,
	}
	mock.lockGetUser.Lock()
	mock.calls.GetUser = append(mock.calls.GetUser, callInfo)
	mock.lockGetUser.Unlock()
	return mock.GetUserFunc(ctx, userId)
}

// GetUserCalls gets all the calls that were made to GetUser.
// Check the length with:
//
//	len(mockedGetUserService.GetUserCalls())
func (mock *GetUserServiceMock) GetUserCalls() []struct {
	Ctx    context.Context
	UserId entity.UserId
} {
	var calls []struct {
		Ctx    context.Context
		UserId entity.UserId
	}
	mock.lockGetUser.RLock()
	calls = mock.calls.GetUser
	mock.lockGetUser.RUnlock()
	return calls
}
