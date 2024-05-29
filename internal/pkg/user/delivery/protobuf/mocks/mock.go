// Code generated by MockGen. DO NOT EDIT.
// Source: auth_grpc.pb.go

// Package mock_grpc is a generated GoMock package.
package mock_grpc

import (
	context "context"
	reflect "reflect"

	grpc "github.com/go-park-mail-ru/2024_1_IMAO/internal/pkg/user/delivery/protobuf"
	gomock "github.com/golang/mock/gomock"
	grpc0 "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockAuthClient is a mock of AuthClient interface.
type MockAuthClient struct {
	ctrl     *gomock.Controller
	recorder *MockAuthClientMockRecorder
}

// MockAuthClientMockRecorder is the mock recorder for MockAuthClient.
type MockAuthClientMockRecorder struct {
	mock *MockAuthClient
}

// NewMockAuthClient creates a new mock instance.
func NewMockAuthClient(ctrl *gomock.Controller) *MockAuthClient {
	mock := &MockAuthClient{ctrl: ctrl}
	mock.recorder = &MockAuthClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthClient) EXPECT() *MockAuthClientMockRecorder {
	return m.recorder
}

// EditEmail mocks base method.
func (m *MockAuthClient) EditEmail(ctx context.Context, in *grpc.EditEmailRequest, opts ...grpc0.CallOption) (*grpc.User, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EditEmail", varargs...)
	ret0, _ := ret[0].(*grpc.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditEmail indicates an expected call of EditEmail.
func (mr *MockAuthClientMockRecorder) EditEmail(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditEmail", reflect.TypeOf((*MockAuthClient)(nil).EditEmail), varargs...)
}

// GetCurrentUser mocks base method.
func (m *MockAuthClient) GetCurrentUser(ctx context.Context, in *grpc.SessionData, opts ...grpc0.CallOption) (*grpc.AuthUser, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetCurrentUser", varargs...)
	ret0, _ := ret[0].(*grpc.AuthUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentUser indicates an expected call of GetCurrentUser.
func (mr *MockAuthClientMockRecorder) GetCurrentUser(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentUser", reflect.TypeOf((*MockAuthClient)(nil).GetCurrentUser), varargs...)
}

// Login mocks base method.
func (m *MockAuthClient) Login(ctx context.Context, in *grpc.ExistedUserData, opts ...grpc0.CallOption) (*grpc.LoggedUser, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Login", varargs...)
	ret0, _ := ret[0].(*grpc.LoggedUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthClientMockRecorder) Login(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthClient)(nil).Login), varargs...)
}

// Logout mocks base method.
func (m *MockAuthClient) Logout(ctx context.Context, in *grpc.SessionData, opts ...grpc0.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Logout", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Logout indicates an expected call of Logout.
func (mr *MockAuthClientMockRecorder) Logout(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockAuthClient)(nil).Logout), varargs...)
}

// Signup mocks base method.
func (m *MockAuthClient) Signup(ctx context.Context, in *grpc.NewUserData, opts ...grpc0.CallOption) (*grpc.LoggedUser, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Signup", varargs...)
	ret0, _ := ret[0].(*grpc.LoggedUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Signup indicates an expected call of Signup.
func (mr *MockAuthClientMockRecorder) Signup(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signup", reflect.TypeOf((*MockAuthClient)(nil).Signup), varargs...)
}

// MockAuthServer is a mock of AuthServer interface.
type MockAuthServer struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServerMockRecorder
}

// MockAuthServerMockRecorder is the mock recorder for MockAuthServer.
type MockAuthServerMockRecorder struct {
	mock *MockAuthServer
}

// NewMockAuthServer creates a new mock instance.
func NewMockAuthServer(ctrl *gomock.Controller) *MockAuthServer {
	mock := &MockAuthServer{ctrl: ctrl}
	mock.recorder = &MockAuthServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthServer) EXPECT() *MockAuthServerMockRecorder {
	return m.recorder
}

// EditEmail mocks base method.
func (m *MockAuthServer) EditEmail(arg0 context.Context, arg1 *grpc.EditEmailRequest) (*grpc.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditEmail", arg0, arg1)
	ret0, _ := ret[0].(*grpc.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditEmail indicates an expected call of EditEmail.
func (mr *MockAuthServerMockRecorder) EditEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditEmail", reflect.TypeOf((*MockAuthServer)(nil).EditEmail), arg0, arg1)
}

// GetCurrentUser mocks base method.
func (m *MockAuthServer) GetCurrentUser(arg0 context.Context, arg1 *grpc.SessionData) (*grpc.AuthUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentUser", arg0, arg1)
	ret0, _ := ret[0].(*grpc.AuthUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentUser indicates an expected call of GetCurrentUser.
func (mr *MockAuthServerMockRecorder) GetCurrentUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentUser", reflect.TypeOf((*MockAuthServer)(nil).GetCurrentUser), arg0, arg1)
}

// Login mocks base method.
func (m *MockAuthServer) Login(arg0 context.Context, arg1 *grpc.ExistedUserData) (*grpc.LoggedUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(*grpc.LoggedUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockAuthServerMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockAuthServer)(nil).Login), arg0, arg1)
}

// Logout mocks base method.
func (m *MockAuthServer) Logout(arg0 context.Context, arg1 *grpc.SessionData) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Logout indicates an expected call of Logout.
func (mr *MockAuthServerMockRecorder) Logout(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockAuthServer)(nil).Logout), arg0, arg1)
}

// Signup mocks base method.
func (m *MockAuthServer) Signup(arg0 context.Context, arg1 *grpc.NewUserData) (*grpc.LoggedUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Signup", arg0, arg1)
	ret0, _ := ret[0].(*grpc.LoggedUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Signup indicates an expected call of Signup.
func (mr *MockAuthServerMockRecorder) Signup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signup", reflect.TypeOf((*MockAuthServer)(nil).Signup), arg0, arg1)
}

// mustEmbedUnimplementedAuthServer mocks base method.
func (m *MockAuthServer) mustEmbedUnimplementedAuthServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedAuthServer")
}

// mustEmbedUnimplementedAuthServer indicates an expected call of mustEmbedUnimplementedAuthServer.
func (mr *MockAuthServerMockRecorder) mustEmbedUnimplementedAuthServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedAuthServer", reflect.TypeOf((*MockAuthServer)(nil).mustEmbedUnimplementedAuthServer))
}

// MockUnsafeAuthServer is a mock of UnsafeAuthServer interface.
type MockUnsafeAuthServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeAuthServerMockRecorder
}

// MockUnsafeAuthServerMockRecorder is the mock recorder for MockUnsafeAuthServer.
type MockUnsafeAuthServerMockRecorder struct {
	mock *MockUnsafeAuthServer
}

// NewMockUnsafeAuthServer creates a new mock instance.
func NewMockUnsafeAuthServer(ctrl *gomock.Controller) *MockUnsafeAuthServer {
	mock := &MockUnsafeAuthServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeAuthServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeAuthServer) EXPECT() *MockUnsafeAuthServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedAuthServer mocks base method.
func (m *MockUnsafeAuthServer) mustEmbedUnimplementedAuthServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedAuthServer")
}

// mustEmbedUnimplementedAuthServer indicates an expected call of mustEmbedUnimplementedAuthServer.
func (mr *MockUnsafeAuthServerMockRecorder) mustEmbedUnimplementedAuthServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedAuthServer", reflect.TypeOf((*MockUnsafeAuthServer)(nil).mustEmbedUnimplementedAuthServer))
}