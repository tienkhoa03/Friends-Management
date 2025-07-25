// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	entity "BE_Friends_Management/internal/domain/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFriendshipService is a mock of FriendshipService interface.
type MockFriendshipService struct {
	ctrl     *gomock.Controller
	recorder *MockFriendshipServiceMockRecorder
}

// MockFriendshipServiceMockRecorder is the mock recorder for MockFriendshipService.
type MockFriendshipServiceMockRecorder struct {
	mock *MockFriendshipService
}

// NewMockFriendshipService creates a new mock instance.
func NewMockFriendshipService(ctrl *gomock.Controller) *MockFriendshipService {
	mock := &MockFriendshipService{ctrl: ctrl}
	mock.recorder = &MockFriendshipServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFriendshipService) EXPECT() *MockFriendshipServiceMockRecorder {
	return m.recorder
}

// CountFriends mocks base method.
func (m *MockFriendshipService) CountFriends(friendsList []*entity.User) int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountFriends", friendsList)
	ret0, _ := ret[0].(int64)
	return ret0
}

// CountFriends indicates an expected call of CountFriends.
func (mr *MockFriendshipServiceMockRecorder) CountFriends(friendsList interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountFriends", reflect.TypeOf((*MockFriendshipService)(nil).CountFriends), friendsList)
}

// CreateFriendship mocks base method.
func (m *MockFriendshipService) CreateFriendship(email1, email2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFriendship", email1, email2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateFriendship indicates an expected call of CreateFriendship.
func (mr *MockFriendshipServiceMockRecorder) CreateFriendship(email1, email2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFriendship", reflect.TypeOf((*MockFriendshipService)(nil).CreateFriendship), email1, email2)
}

// RetrieveCommonFriends mocks base method.
func (m *MockFriendshipService) RetrieveCommonFriends(email1, email2 string) ([]*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveCommonFriends", email1, email2)
	ret0, _ := ret[0].([]*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveCommonFriends indicates an expected call of RetrieveCommonFriends.
func (mr *MockFriendshipServiceMockRecorder) RetrieveCommonFriends(email1, email2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveCommonFriends", reflect.TypeOf((*MockFriendshipService)(nil).RetrieveCommonFriends), email1, email2)
}

// RetrieveFriendsList mocks base method.
func (m *MockFriendshipService) RetrieveFriendsList(email string) ([]*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetrieveFriendsList", email)
	ret0, _ := ret[0].([]*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RetrieveFriendsList indicates an expected call of RetrieveFriendsList.
func (mr *MockFriendshipServiceMockRecorder) RetrieveFriendsList(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetrieveFriendsList", reflect.TypeOf((*MockFriendshipService)(nil).RetrieveFriendsList), email)
}
