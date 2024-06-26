// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mock_usecases is a generated GoMock package.
package mock_usecases

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2024_1_IMAO/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockUsersStorageInterface is a mock of UsersStorageInterface interface.
type MockUsersStorageInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUsersStorageInterfaceMockRecorder
}

// MockUsersStorageInterfaceMockRecorder is the mock recorder for MockUsersStorageInterface.
type MockUsersStorageInterfaceMockRecorder struct {
	mock *MockUsersStorageInterface
}

// NewMockUsersStorageInterface creates a new mock instance.
func NewMockUsersStorageInterface(ctrl *gomock.Controller) *MockUsersStorageInterface {
	mock := &MockUsersStorageInterface{ctrl: ctrl}
	mock.recorder = &MockUsersStorageInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsersStorageInterface) EXPECT() *MockUsersStorageInterfaceMockRecorder {
	return m.recorder
}

// AddSession mocks base method.
func (m *MockUsersStorageInterface) AddSession(ctx context.Context, id uint) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSession", ctx, id)
	ret0, _ := ret[0].(string)
	return ret0
}

// AddSession indicates an expected call of AddSession.
func (mr *MockUsersStorageInterfaceMockRecorder) AddSession(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSession", reflect.TypeOf((*MockUsersStorageInterface)(nil).AddSession), ctx, id)
}

// CreateUser mocks base method.
func (m *MockUsersStorageInterface) CreateUser(ctx context.Context, email, password, passwordRepeat string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, email, password, passwordRepeat)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUsersStorageInterfaceMockRecorder) CreateUser(ctx, email, password, passwordRepeat interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUsersStorageInterface)(nil).CreateUser), ctx, email, password, passwordRepeat)
}

// EditUserEmail mocks base method.
func (m *MockUsersStorageInterface) EditUserEmail(ctx context.Context, id uint, email string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditUserEmail", ctx, id, email)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditUserEmail indicates an expected call of EditUserEmail.
func (mr *MockUsersStorageInterfaceMockRecorder) EditUserEmail(ctx, id, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditUserEmail", reflect.TypeOf((*MockUsersStorageInterface)(nil).EditUserEmail), ctx, id, email)
}

// GetLastID mocks base method.
func (m *MockUsersStorageInterface) GetLastID(ctx context.Context) uint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastID", ctx)
	ret0, _ := ret[0].(uint)
	return ret0
}

// GetLastID indicates an expected call of GetLastID.
func (mr *MockUsersStorageInterfaceMockRecorder) GetLastID(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastID", reflect.TypeOf((*MockUsersStorageInterface)(nil).GetLastID), ctx)
}

// GetUserByEmail mocks base method.
func (m *MockUsersStorageInterface) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUsersStorageInterfaceMockRecorder) GetUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUsersStorageInterface)(nil).GetUserByEmail), ctx, email)
}

// GetUserBySession mocks base method.
func (m *MockUsersStorageInterface) GetUserBySession(ctx context.Context, sessionID string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBySession", ctx, sessionID)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBySession indicates an expected call of GetUserBySession.
func (mr *MockUsersStorageInterfaceMockRecorder) GetUserBySession(ctx, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBySession", reflect.TypeOf((*MockUsersStorageInterface)(nil).GetUserBySession), ctx, sessionID)
}

// RemoveSession mocks base method.
func (m *MockUsersStorageInterface) RemoveSession(ctx context.Context, sessionID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSession", ctx, sessionID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveSession indicates an expected call of RemoveSession.
func (mr *MockUsersStorageInterfaceMockRecorder) RemoveSession(ctx, sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSession", reflect.TypeOf((*MockUsersStorageInterface)(nil).RemoveSession), ctx, sessionID)
}

// SessionExists mocks base method.
func (m *MockUsersStorageInterface) SessionExists(sessionID string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SessionExists", sessionID)
	ret0, _ := ret[0].(bool)
	return ret0
}

// SessionExists indicates an expected call of SessionExists.
func (mr *MockUsersStorageInterfaceMockRecorder) SessionExists(sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SessionExists", reflect.TypeOf((*MockUsersStorageInterface)(nil).SessionExists), sessionID)
}

// UserExists mocks base method.
func (m *MockUsersStorageInterface) UserExists(ctx context.Context, email string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserExists", ctx, email)
	ret0, _ := ret[0].(bool)
	return ret0
}

// UserExists indicates an expected call of UserExists.
func (mr *MockUsersStorageInterfaceMockRecorder) UserExists(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserExists", reflect.TypeOf((*MockUsersStorageInterface)(nil).UserExists), ctx, email)
}
